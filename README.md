# ciel
A tool for controlling multi-layer file systems and containers

## manual

During the rapid iteration (before version 1.x.x), you may read help messages:
```
ciel help
```

Or Wiki:
- https://github.com/AOSC-Dev/ciel/wiki/The-Ciel-User-Manual-en
- https://github.com/AOSC-Dev/ciel/wiki/The-Ciel-User-Manual-zh_CN

## installation
```
make
sudo make install
```

## dependencies

Building:
- git
- make
- Go

Runtime:
- systemd's container components
- overlayfs (kernel module)
- coreutils
- tar
