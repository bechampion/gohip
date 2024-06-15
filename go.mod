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

require github.com/bechampion/gohip/types v0.0.0-20240615162333-918ade08d726
