Documentation: http://bgb.bircd.org/pandocs.htm

Profiling:

    PROFILE=$(pwd) $(make run-gbc-cmd) data/bomberman.gb
    go tool pprof -http localhost:8081 cpu.pprof
