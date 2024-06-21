package timescale

import (
	"context"
	"fmt"

	_ "embed"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tannerryan/davisweather"
)

type TimescaleClient struct {
	Pool *pgxpool.Pool
}

func NewTimescaleClient(ctx context.Context, connString string) (*TimescaleClient, error) {
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	cfg.ConnConfig.Tracer = otelpgx.NewTracer()

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &TimescaleClient{Pool: pool}, nil
}

func (c *TimescaleClient) Close() {
	c.Pool.Close()
}

//go:embed queries/insert_vp2p.pgsql
var insertVantagePro2Plus string

// Report is the latest weather report.
// type Report struct {
// 	DeviceID  string    `json:"deviceID"`  // DeviceID is unique device ID
// 	Timestamp time.Time `json:"timestamp"` // Timestamp is the time the Report was last modified

// 	Temperature *float64 `json:"temperature"` // Temperature (°F)
// 	Humidity    *float64 `json:"humidity"`    // Humidity (%RH)
// 	Dewpoint    *float64 `json:"dewpoint"`    // Dewpoint (°F)
// 	Wetbulb     *float64 `json:"wetbulb"`     // Wetbulb (°F)
// 	HeatIndex   *float64 `json:"heatindex"`   // HeatIndex (°F)
// 	WindChill   *float64 `json:"windchill"`   // WindChill (°F)
// 	THWIndex    *float64 `json:"thwIndex"`    // THWIndex is "feels like" (°F)
// 	THSWIndex   *float64 `json:"thswIndex"`   // THWSIndex is "feels like" including solar (°F)

// 	WindSpeedLast          *float64 `json:"windSpeedLast"`          // WindSpeedLast is most recent wind speed (mph)
// 	WindDirLast            *float64 `json:"windDirLast"`            // WindDirLast is most recent wind direction (°)
// 	WindSpeedAvgLast1Min   *float64 `json:"windSpeedAvg1Min"`       // WindSpeedAvgLast1Min is average wind over last minute (mph)
// 	WindDirAvgLast1Min     *float64 `json:"windDirAvg1Min"`         // WindDirAvgLast1Min is average wind direction over last minute (°)
// 	WindSpeedAvgLast2Min   *float64 `json:"windSpeedAvg2Min"`       // WindSpeedAvgLast2Min is average wind over last 2 minutes (mph)
// 	WindDirAvgLast2Min     *float64 `json:"windDirAvg2Min"`         // WindDirAvgLast2Min is average wind direction over last 2 minutes (°)
// 	WindSpeedHighLast2Min  *float64 `json:"windGustSpeedLast2Min"`  // WindSpeedHighLast2Min is max gust over last 2 minutes (mph)
// 	WindDirAtHighLast2Min  *float64 `json:"windGustDirLast2Min"`    // WindDirAtHighLast2Min is max gust direction over last 2 minutes (°)
// 	WindSpeedAvgLast10Min  *float64 `json:"windSpeedAvg10Min"`      // WindSpeedAvgLast10Min is average wind over last 10 minutes (mph)
// 	WindDirAvgLast10Min    *float64 `json:"windDirAvg10Min"`        // WindDirAvgLast10Min is average wind dir over last 10 minutes (°)
// 	WindSpeedHighLast10Min *float64 `json:"windGustSpeedLast10Min"` // WindSpeedHighLast10Min is max gust over last 10 minutes (mph)
// 	WindDirAtHighLast10Min *float64 `json:"windGustDirLast10Min"`   // WindDirAtHighLast10Min is max gust direction over last 10 minutes (°)

// 	RainSize              *float64   `json:"rainSize"`              // RainSize is size of rain collector (1: 0.01", 2: 0.2mm)
// 	RainRateLast          *float64   `json:"rainRateLast"`          // RainRateLast is most recent rain rate (count/hour)
// 	RainRateHigh          *float64   `json:"rainRateHigh"`          // RainRateHigh is highest rain rate over last minute (count/hour)
// 	RainLast15Min         *float64   `json:"rainLast15Min"`         // RainLast15Min is rain count in last 15 minutes (count)
// 	RainRateHighLast15Min *float64   `json:"rainRateHighLast15Min"` // RainRateHighLast15Min is highest rain count rate over last 15 minutes (count/hour)
// 	RainLast60Min         *float64   `json:"rainLast60Min"`         // RainLast60Min is rain count over last 60 minutes (count)
// 	RainLast24Hour        *float64   `json:"rainLast24Hour"`        // RainLast24Hour is rain count over last 24 hours (count)
// 	RainStorm             *float64   `json:"rainStorm"`             // RainStorm is rain since last 24 hour break in rain (count)
// 	RainStormStartAt      *time.Time `json:"rainStormStart"`        // RainStormStartAt is time of rain storm start

// 	SolarRad *float64 `json:"solarRad"` // SolarRad is solar radiation (W/m²)
// 	UVIndex  *float64 `json:"uvIndex"`  // UVIndex is solar UV index

// 	RXState          string `json:"signal"`  // RXState is ISS receiver status
// 	TransBatteryFlag string `json:"battery"` // TransBatteryFlag is ISS battery status

// 	RainfallDaily        *float64   `json:"rainDaily"`          // RainfallDaily is total rain since midnight (count)
// 	RainfallMonthly      *float64   `json:"rainMonthly"`        // RainfallMonthly is total rain since first of month (count)
// 	RainfallYear         *float64   `json:"rainYear"`           // RainfallYear is total rain since first of year (count)
// 	RainStormLast        *float64   `json:"rainStormLast"`      // RainStormLast is rain since last 24 hour break in rain (count)
// 	RainStormLastStartAt *time.Time `json:"rainStormLastStart"` // RainStormLastStartAt is time of last rain storm start
// 	RainStormLastEndAt   *time.Time `json:"rainStormLastEnd"`   // rainStormLastEndAt is time of last rain storm end

// 	BarometerSeaLevel *float64 `json:"barometerSeaLevel"` // BarometerSeaLevel is barometer reading with elevation adjustment (inches)
// 	BarometerTrend    *float64 `json:"barometerTrend"`    // BarometerTrend is 3 hour barometric trend (inches)
// 	BarometerAbsolute *float64 `json:"barometerAbsolute"` // BarometerAbsolute is barometer reading at current elevation (inches)

// 	TemperatureIndoor *float64 `json:"indoorTemperature"` // TemperatureIndoor is indoor temp (°F)
// 	HumidityIndoor    *float64 `json:"indoorHumidity"`    // HumidityIndoor is indoor humidity (%)
// 	DewPointIndoor    *float64 `json:"indoorDewpoint"`    // DewPointIndoor is indoor dewpoint (°F)
// 	HeatIndexIndoor   *float64 `json:"indoorHeatIndex"`   // HeatIndexIndoor is indoor heat index (°F)

// 	notify       chan bool   // notify emits a boolean when the Report contents are modified
// 	verbose      bool        // verbose enables Report logging to stdout
// 	lastChecksum string      // lastChecksum is MD5 checksum of the Report state
// 	lastBytes    []byte      // lastBytes is the JSON representation of the Report state
// 	mutex        *sync.Mutex // mutex is for atomic report actions
// }

func (c *TimescaleClient) Insert(ctx context.Context, report *davisweather.Report) error {
	_, err := c.Pool.Exec(ctx, insertVantagePro2Plus,
		report.Timestamp,
		report.DeviceID,
		report.Temperature,
		report.Humidity,
		report.Dewpoint,
		report.Wetbulb,
		report.HeatIndex,
		report.WindChill,
		report.THWIndex,
		report.THSWIndex,
		report.WindSpeedLast,
		report.WindDirLast,
		report.WindSpeedAvgLast1Min,
		report.WindDirAvgLast1Min,
		report.WindSpeedAvgLast2Min,
		report.WindDirAvgLast2Min,
		report.WindSpeedHighLast2Min,
		report.WindDirAtHighLast2Min,
		report.WindSpeedAvgLast10Min,
		report.WindDirAvgLast10Min,
		report.WindSpeedHighLast10Min,
		report.WindDirAtHighLast10Min,
		report.RainSize,
		report.RainRateLast,
		report.RainRateHigh,
		report.RainLast15Min,
		report.RainRateHighLast15Min,
		report.RainLast60Min,
		report.RainLast24Hour,
		report.RainStorm,
		report.RainStormStartAt,
		report.SolarRad,
		report.UVIndex,
		report.RXState,
		report.TransBatteryFlag,
		report.RainfallDaily,
		report.RainfallMonthly,
		report.RainfallYear,
		report.RainStormLast,
		report.RainStormLastStartAt,
		report.RainStormLastEndAt,
		report.BarometerSeaLevel,
		report.BarometerTrend,
		report.BarometerAbsolute,
		report.TemperatureIndoor,
		report.HumidityIndoor,
		report.DewPointIndoor,
		report.HeatIndexIndoor,
	)
	if err != nil {
		return fmt.Errorf("insert data: %w", err)
	}

	return nil
}
