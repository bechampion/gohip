module hip

go 1.22.3

// replace osdata => ./osdata

// replace systemd => ./systemd

// replace types => ./types

// replace others => ./others

// require (
// 	osdata v0.0.0-00010101000000-000000000000
// 	others v0.0.0-00010101000000-000000000000
// 	systemd v0.0.0-00010101000000-000000000000
// 	types v0.0.0-00010101000000-000000000000
// )

require (
	github.com/bechampion/gohip/osdata v0.0.0-20240615162838-24fd4fe6ae1c
	github.com/bechampion/gohip/systemd v0.0.0-20240615162649-e52f38b1aa94
	github.com/bechampion/gohip/types v0.0.0-20240615162649-e52f38b1aa94
)

require (
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/power-devops/perfstat v0.0.0-20210106213030-5aafc221ea8c // indirect
	github.com/shirou/gopsutil/v3 v3.24.4 // indirect
	github.com/shoenig/go-m1cpu v0.1.6 // indirect
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	golang.org/x/sys v0.19.0 // indirect
)
