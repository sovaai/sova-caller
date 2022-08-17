# GO Bindings

1. [Download latest PJSIP](https://www.pjsip.org/download.htm)
2. Extract to ${pjprject-src-dir}
3. [Make and Install PJSIP](https://trac.pjsip.org/repos/wiki/Getting-Started)
4. mkdir pjsua2
5. cd pjsua2
6. cp ${pjprject-src-dir}/pjsip-apps/src/swig/pjsua2.i .
7. cp ${pjprject-src-dir}/pjsip-apps/src/swig/symbols.i .
8. Generate go bindings

```
swig -c++ -go -cgo -intgosize 64 \
   -outcurrentdir \
   -I/usr/local/include \
   pjsua2.i
```

9. Add static libs - linux

```
package pjsua2

/*
#cgo CPPFLAGS: -DPJ_AUTOCONF=1 -O2 -DPJ_IS_BIG_ENDIAN=0
#cgo CPPFLAGS: -DPJ_IS_LITTLE_ENDIAN=1
#cgo CPPFLAGS: -DPJMEDIA_USE_OLD_FFMPEG=1
#cgo CPPFLAGS: -I/usr/local/include -I/usr/include
#cgo CPPFLAGS: -I/usr/include/c++/5 -I/usr/include/x86_64-linux-gnu/c++/5
#cgo LDFLAGS: -L/usr/local/lib
#cgo LDFLAGS: -lpjsua2-x86_64-unknown-linux-gnu
#cgo LDFLAGS: -lpjsua-x86_64-unknown-linux-gnu
#cgo LDFLAGS: -lpjsip-ua-x86_64-unknown-linux-gnu
#cgo LDFLAGS: -lpjsip-simple-x86_64-unknown-linux-gnu
#cgo LDFLAGS: -lpjsip-x86_64-unknown-linux-gnu
#cgo LDFLAGS: -lpjmedia-codec-x86_64-unknown-linux-gnu
#cgo LDFLAGS: -lpjmedia-videodev-x86_64-unknown-linux-gnu
#cgo LDFLAGS: -lpjmedia-audiodev-x86_64-unknown-linux-gnu
#cgo LDFLAGS: -lpjmedia-x86_64-unknown-linux-gnu
#cgo LDFLAGS: -lpjnath-x86_64-unknown-linux-gnu
#cgo LDFLAGS: -lpjlib-util-x86_64-unknown-linux-gnu
#cgo LDFLAGS: -lsrtp-x86_64-unknown-linux-gnu
#cgo LDFLAGS: -lresample-x86_64-unknown-linux-gnu
#cgo LDFLAGS: -lgsmcodec-x86_64-unknown-linux-gnu
#cgo LDFLAGS: -lspeex-x86_64-unknown-linux-gnu
#cgo LDFLAGS: -lilbccodec-x86_64-unknown-linux-gnu
#cgo LDFLAGS: -lg7221codec-x86_64-unknown-linux-gnu
#cgo LDFLAGS: -lyuv-x86_64-unknown-linux-gnu
#cgo LDFLAGS: -lwebrtc-x86_64-unknown-linux-gnu
#cgo LDFLAGS: -lpj-x86_64-unknown-linux-gnu
#cgo LDFLAGS: -lssl -lcrypto -luuid -lm -lrt -lpthread -lasound
```

10. Install missing libs (ex. uuid-dev,libasound2-dev)
11. go build
