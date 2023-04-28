package sonus

import (
	"context"
	"os"
	"testing"
)

var testSBC *SBC

func init() {
	testSBC = NewSBC(os.Getenv("SONUS_TARGET"), os.Getenv("SONUS_USER"), os.Getenv("SONUS_PASSWORD"))
}

func TestNewSBC(t *testing.T) {
	type args struct {
		address  string
		user     string
		password string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test new SBC",
			args: args{
				address:  os.Getenv("SONUS_TARGET"),
				user:     os.Getenv("SONUS_USER"),
				password: os.Getenv("SONUS_PASSWORD"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sbc := NewSBC(tt.args.address, tt.args.user, tt.args.password)
			if sbc == nil {
				t.Errorf("NewSBC() = %v", sbc)
			}
			if sbc.AddressContexts == nil || len(sbc.AddressContexts.AddressContext) == 0 {
				t.Errorf("NewSBC() no contexts found: %v", sbc.AddressContexts)
			}
		})
	}
}

func TestZoneStatus(t *testing.T) {
	for _, aCtx := range testSBC.AddressContexts.AddressContext {
		stats := &ZoneStats{}
		err := testSBC.GetAndParse(context.Background(), stats, zoneStatusPath, aCtx.Name)
		if err != nil {
			t.Errorf("TestZoneStatus returned error: %v", err)
		}
		// for _, z := range stats.Zone {
		// 	for _, sipStats := range z.SipCurrentStatistics {
		// 		processStruct(sipStats)
		// 	}
		// }
	}
}
