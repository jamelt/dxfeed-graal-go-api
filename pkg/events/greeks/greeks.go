package greeks

import (
	"math"
	"strconv"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/formatutil"
)

type Greeks struct {
	eventSymbol *string
	eventTime   int64
	eventFlags  int32
	index       int64
	price       float64
	volatility  float64
	delta       float64
	gamma       float64
	theta       float64
	rho         float64
	vega        float64
}

func NewGreeks(eventSymbol string) *Greeks {
	return &Greeks{
		eventSymbol: &eventSymbol,
		price:       math.NaN(),
		volatility:  math.NaN(),
		delta:       math.NaN(),
		gamma:       math.NaN(),
		theta:       math.NaN(),
		rho:         math.NaN(),
		vega:        math.NaN(),
	}
}

func (g *Greeks) Type() eventcodes.EventCode {
	return eventcodes.Greeks
}

func (g *Greeks) EventSymbol() *string {
	return g.eventSymbol
}

func (g *Greeks) SetEventSymbol(eventSymbol string) {
	*g.eventSymbol = eventSymbol
}

func (g *Greeks) EventTime() int64 {
	return g.eventTime
}

func (g *Greeks) SetEventTime(eventTime int64) {
	g.eventTime = eventTime
}

func (g *Greeks) EventFlags() int32 {
	return g.eventFlags
}

func (g *Greeks) SetEventFlags(eventFlags int32) {
	g.eventFlags = eventFlags
}

func (g *Greeks) Index() int64 {
	return g.index
}

func (g *Greeks) SetIndex(index int64) {
	g.index = index
}

func (g *Greeks) Price() float64 {
	return g.price
}

func (g *Greeks) SetPrice(price float64) {
	g.price = price
}

func (g *Greeks) Volatility() float64 {
	return g.volatility
}

func (g *Greeks) SetVolatility(volatility float64) {
	g.volatility = volatility
}

func (g *Greeks) Delta() float64 {
	return g.delta
}

func (g *Greeks) SetDelta(delta float64) {
	g.delta = delta
}

func (g *Greeks) Gamma() float64 {
	return g.gamma
}

func (g *Greeks) SetGamma(gamma float64) {
	g.gamma = gamma
}

func (g *Greeks) Theta() float64 {
	return g.theta
}

func (g *Greeks) SetTheta(theta float64) {
	g.theta = theta
}

func (g *Greeks) Rho() float64 {
	return g.rho
}

func (g *Greeks) SetRho(rho float64) {
	g.rho = rho
}

func (g *Greeks) Vega() float64 {
	return g.vega
}

func (g *Greeks) SetVega(vega float64) {
	g.vega = vega
}

func (g *Greeks) String() string {
	return "Greeks{" + formatutil.FormatString(g.EventSymbol()) +
		", eventTime=" + formatutil.FormatTime(g.EventTime()) +
		", eventFlags=" + strconv.FormatInt(int64(g.eventFlags), 10) +
		", index=" + strconv.FormatInt(g.index, 10) +
		", price=" + formatutil.FormatFloat64(g.price) +
		", volatility=" + formatutil.FormatFloat64(g.volatility) +
		", delta=" + formatutil.FormatFloat64(g.delta) +
		", gamma=" + formatutil.FormatFloat64(g.gamma) +
		", theta=" + formatutil.FormatFloat64(g.theta) +
		", rho=" + formatutil.FormatFloat64(g.rho) +
		", vega=" + formatutil.FormatFloat64(g.vega) +
		"}"
}
