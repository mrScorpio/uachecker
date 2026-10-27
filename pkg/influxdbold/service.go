package influxdbold

import (
	"time"

	oldb "github.com/influxdata/influxdb/client/v2"
)

type InfluxCl struct {
	Cl   oldb.HTTPClient
	Bp   oldb.BatchPoints
	Tags map[string]string
}

func NewDbConn(addr, user, pass string) (*InfluxCl, error) {
	cl, err := oldb.NewHTTPClient(oldb.HTTPConfig{
		Addr:     addr,
		Username: user,
		Password: pass,
	})
	if err != nil {
		return nil, err
	}
	bp, err := oldb.NewBatchPoints(oldb.BatchPointsConfig{
		Database:  "stendb",
		Precision: "ms",
	})
	if err != nil {
		return nil, err
	}

	return &InfluxCl{Cl: cl, Bp: bp, Tags: map[string]string{"sys": "gtd"}}, nil
}

func (iCl *InfluxCl) AddPoint(p map[string]any) error {
	pt, err := oldb.NewPoint("values", iCl.Tags, p, time.Now())
	if err != nil {
		return err
	}
	iCl.Bp.AddPoint(pt)
	if len(iCl.Bp.Points()) > 6 {
		err := iCl.Cl.Write(iCl.Bp)
		if err != nil {
			return err
		}
		iCl.Bp, err = oldb.NewBatchPoints(oldb.BatchPointsConfig{
			Database:  "stendb",
			Precision: "ms",
		})
		if err != nil {
			return err
		}
	}
	return nil
}
