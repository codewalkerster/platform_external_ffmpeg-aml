package ffmpeg

import (
    "android/soong/android"
    "android/soong/cc"
    "os/exec"
    "fmt"
    "strings"
)

func init() {
    android.RegisterModuleType("ffmpeg_defaults", ffmpegDefaultsFactory)
}

func ffmpegDefaultsFactory() (android.Module) {
    module := cc.DefaultsFactory()
    android.AddLoadHook(module, ffmpegDefaults)
    return module
}

func ffmpegDefaults(ctx android.LoadHookContext) {
    type props struct {
        Cflags []string
        //Srcs []string
    }
    p := &props{}
    p.Cflags = getCflags(ctx)
    //p.Srcs = getSrcs(ctx)
    ctx.AppendProperties(p)
}

func getCflags(ctx android.BaseContext) ([]string) {
    var cflags []string
    var arch = string(ctx.DeviceConfig().DeviceArch())
    if arch == "arm" || arch == "arm64" {
        cflags = append(cflags,
                        "-std=c99",
                        "-llvm",
                        "-DDISABLE_NEONINTR",
                        "-D_FILE_OFFSET_BITS=64",
                        "-D_LARGEFILE_SOURCE",
                        "-D_ISOC99_SOURCE",
                        "-D_POSIX_C_SOURCE=200112")
    }

    var dir = string(ctx.Config().Getenv("PWD"))
    gcmd := "cd " + dir + "/" + ctx.ModuleDir() + "&& git log | grep commit -m 1 | cut -d' ' -f 2"
    gitout, giterr := exec.Command("/bin/bash", "-c", gcmd).CombinedOutput()
    if giterr != nil {
        fmt.Printf("get git err, gitout %s\n", gitout)
    } else {
        // No build date here: it changes on every soong run, which changes
        // the clang command line and makes ninja rebuild the whole module.
        // The git hash is stable across builds, so it is kept.
        gitversion := strings.TrimSpace(string(gitout))
        buildname := strings.TrimSpace(string(ctx.Config().Getenv("LOGNAME")))
        ver := "-DGIT_INFO=" + "\"" + gitversion + " (" + buildname + ")" + "\""
        cflags = append(cflags, ver)
    }

    return cflags
}
