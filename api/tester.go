package api

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/CybexDex/cybex-go/types"
	"github.com/denkhaus/logging"
	sort "github.com/emirpasic/gods/utils"
	deadlock "github.com/sasha-s/go-deadlock"
	"gopkg.in/tomb.v2"
)

var (
	//LoopSeconds = time for one pass to calc dynamic delay
	LoopSeconds = 60
	//our known node endpoints
	knownEndpoints = []string{
		"wss://eu.openledger.info/ws",
		"wss://cybex.openledger.info/ws",
		"wss://dex.rnglab.org",
		"wss://api.cybex.bhuz.info/ws",
		"wss://cybex.crypto.fans/ws",
		"wss://node.market.rudex.org",
		"wss://api.cyb.blckchnd.com",
		"wss://eu.nodes.cybex.ws",
		"wss://cybws.roelandp.nl/ws",
		"wss://cybfullnode.bangzi.info/ws",
		"wss://api-ru.cyb.blckchnd.com",
		"wss://kc-us-dex.xeldal.com/ws",
		"wss://api.cybxchng.com",
		"wss://api.cyb.network",
		"wss://dexnode.net/ws",
		"wss://us.nodes.cybex.ws",
		"wss://api.cyb.mobi/ws",
		"wss://blockzms.xyz/ws",
		"wss://cyb-api.lafona.net/ws",
		"wss://api.cyb.ai/",
		"wss://la.dexnode.net/ws",
		"wss://openledger.hk/ws",
		"wss://sg.nodes.cybex.ws",
		"wss://cyb.open.icowallet.net/ws",
		"wss://ws.gdex.io",
		"wss://cybex-api.wancloud.io/ws",
		"wss://ws.hellocyb.com/",
		"wss://cybex.dacplay.org/ws",
		"wss://crazybit.online",
		"wss://kimziv.com/ws",
		"wss://wss.ioex.top",
		"wss://node.cybcharts.com/ws",
		"wss://cyb-seoul.clockwork.gr/",
		"wss://cybex.cyberit.io/",
		"wss://api.cybgo.net/ws",
		"wss://ws.winex.pro",
		"wss://cyb.to0l.cn:4443/ws",
		"wss://cybex.cyb123.cc:15138/",
		"wss://bit.cybabc.org/ws",
		"wss://ws.gdex.top",
	}
)

type LatencyTester interface {
	Start()
	Close() error
	String() string
	AddEndpoint(ep string)
	OnTopNodeChanged(fn func(string) error)
	TopNodeEndpoint() string
	TopNodeClient() WebsocketClient
	Done() <-chan struct{}
}

//NodeStats holds stat data for each endpoint
type NodeStats struct {
	cli        WebsocketClient
	latency    time.Duration
	requiredDB []string
	attempts   int64
	errors     int64
	endpoint   string
}

func (p *NodeStats) onError(err error) {
	p.errors++
}

//Latency returns the nodes latency
func (p *NodeStats) Latency() time.Duration {
	if p.attempts > 0 {
		return time.Duration(int64(p.latency) / p.attempts)
	}

	return 0
}

//Score returns reliability score for each node. The less the better.
func (p *NodeStats) Score() int64 {
	lat := int64(p.Latency())

	if lat == 0 {
		return math.MaxInt64
	}

	if p.errors == 0 {
		return lat
	}

	// add 50ms per error
	return lat + p.errors*50000000
}

// String returns the stats string representation
func (p *NodeStats) String() string {
	return fmt.Sprintf("ep: %s | attempts: %d | errors: %d | latency: %s | score: %d",
		p.endpoint, p.attempts, p.errors, p.Latency(), p.Score())
}

//NewNodeStats creates a new stat object
func NewNodeStats(wsRPCEndpoint string) *NodeStats {
	stats := &NodeStats{
		endpoint: wsRPCEndpoint,
		cli:      NewWebsocketClient(wsRPCEndpoint),
		requiredDB: []string{
			"database",
			"history",
			"network_broadcast",
			"crypto",
		},
	}

	stats.cli.OnError(stats.onError)
	return stats
}

func (p *NodeStats) Equals(n *NodeStats) bool {
	return p.endpoint == n.endpoint
}

func (p *NodeStats) checkNode() {
	_, err := p.cli.CallAPI(1, "login", "", "")
	if err != nil {
		p.errors++
	}

	for _, name := range p.requiredDB {
		_, err := p.cli.CallAPI(1, name, types.EmptyParams)
		if err != nil {
			p.errors++
		}
	}
}

func (p *NodeStats) check() {
	p.attempts++
	if err := p.cli.Connect(); err != nil {
		p.errors++
	}
	defer p.cli.Close()

	tm := time.Now()
	p.checkNode()
	p.latency += time.Since(tm)
}

type latencyTester struct {
	mu               deadlock.RWMutex
	tmb              *tomb.Tomb
	toApply          []string
	fallbackURL      string
	onTopNodeChanged func(string) error
	stats            []interface{}
	pass             int
}

func NewLatencyTester(fallbackURL string) (LatencyTester, error) {
	return NewLatencyTesterWithContext(context.Background(), fallbackURL)
}

func NewLatencyTesterWithContext(ctx context.Context, fallbackURL string) (LatencyTester, error) {
	tmb, _ := tomb.WithContext(ctx)
	lat := latencyTester{
		fallbackURL: fallbackURL,
		stats:       make([]interface{}, 0, len(knownEndpoints)),
		tmb:         tmb,
	}

	lat.createStats(knownEndpoints)
	return &lat, nil
}

func (p *latencyTester) String() string {
	builder := strings.Builder{}

	p.mu.Lock()
	defer p.mu.Unlock()
	for _, s := range p.stats {
		stat := s.(*NodeStats)
		builder.WriteString(stat.String())
		builder.WriteString("\n")
	}

	return builder.String()
}

func (p *latencyTester) OnTopNodeChanged(fn func(string) error) {
	p.onTopNodeChanged = fn
}

//AddEndpoint adds a new endpoint while the latencyTester is running
func (p *latencyTester) AddEndpoint(ep string) {
	p.toApply = append(p.toApply, ep)
}

func (p *latencyTester) sortResults(notify bool) error {

	p.mu.Lock()
	oldTop := p.stats[0].(*NodeStats)
	sort.Sort(p.stats, func(a, b interface{}) int {
		sa := a.(*NodeStats).Score()
		sb := b.(*NodeStats).Score()
		if sa > sb {
			return 1
		}

		if sa < sb {
			return -1
		}

		return 0
	})

	newTop := p.stats[0].(*NodeStats)
	p.mu.Unlock()

	if notify && !oldTop.Equals(newTop) {
		if p.onTopNodeChanged != nil {
			return p.onTopNodeChanged(newTop.endpoint)
		}
	}

	return nil
}

func (p *latencyTester) createStats(eps []string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, ep := range eps {
		found := false
		for _, s := range p.stats {
			stat := s.(*NodeStats)
			if stat.endpoint == ep {
				found = true
			}
		}

		if !found {
			p.stats = append(
				p.stats,
				NewNodeStats(ep),
			)
		}
	}
}

//TopNodeEndpoint returns the fastest endpoint URL. If the tester has no validated results
//your given fallback endpoint is returned.
func (p *latencyTester) TopNodeEndpoint() string {
	if p.pass > 0 {
		p.mu.RLock()
		defer p.mu.RUnlock()
		st := p.stats[0].(*NodeStats)
		return st.endpoint
	}

	return p.fallbackURL
}

//TopNodeClient returns a new WebsocketClient to connect to the fastest node.
//If the tester has no validated results, a client with your given
//fallback endpoint is returned. You need to call Connect for yourself.
func (p *latencyTester) TopNodeClient() WebsocketClient {
	return NewWebsocketClient(
		p.TopNodeEndpoint(),
	)
}

// Done returns the channel that can be used to wait until
// the tester has finished.
func (p *latencyTester) Done() <-chan struct{} {
	return p.tmb.Dead()
}

func (p *latencyTester) runPass() error {
	// dynamic sleep time
	slp := time.Duration(LoopSeconds/len(p.stats)) * time.Second
	for i := 0; i < len(p.stats); i++ {
		select {
		case <-p.tmb.Dying():
			return tomb.ErrDying
		default:
			time.Sleep(slp)
			p.mu.RLock()
			st := p.stats[i].(*NodeStats)
			st.check()
			p.mu.RUnlock()
		}
	}

	return nil
}

//Start starts the testing process
func (p *latencyTester) Start() {
	logging.Debug("latencytester: start")

	p.tmb.Go(func() error {
		for {
			select {
			case <-p.tmb.Dying():
				p.sortResults(false)
				return tomb.ErrDying
			default:
				//apply later incoming endpoints
				p.createStats(p.toApply)
				if err := p.runPass(); err != nil {
					//provide sorted results on return
					p.sortResults(false)
					return err
				}

				p.sortResults(true)
				p.pass++
			}
		}
	})
}

//Close stops the tester and waits until all goroutines have finished.
func (p *latencyTester) Close() error {
	logging.Debug("latencytester: kill [tomb]")
	p.tmb.Kill(nil)

	logging.Debug("latencytester: wait [tomb]")
	return p.tmb.Wait()
}
