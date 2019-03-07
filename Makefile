bin := queuesim

FLAGS :=
debug: FLAGS := -gcflags="-N -l"

$(bin): main.go
	go build $(FLAGS) -o $(bin)

debug: $(bin)

.PHONY: clean
clean:
	rm -f $(bin)
