const std = @import("std");
//const zcc = @import("compile_commands");

pub fn build(b: *std.Build) void {
    const target = b.standardTargetOptions(.{});
    var targets = std.ArrayList(*std.Build.Step.Compile).init(b.allocator);

    const optimize = b.standardOptimizeOption(.{});
    // Create executable
    const exe = b.addExecutable(.{
        .name = "urfp",
        .target = target,
        .optimize = optimize,
    });

    // Add C source files
    exe.addCSourceFiles(.{
        .files = &.{
            "src/main.c",
        },
        .flags = &.{"-Wall"},
    });
    exe.linkLibC();
    exe.linkSystemLibrary2("libfprint-2", .{ .use_pkg_config = .force });
    exe.linkSystemLibrary2("glib-2.0", .{ .use_pkg_config = .force });

    targets.append(exe) catch @panic("OOM");

    const install_exe = b.addInstallArtifact(exe, .{});
    b.getInstallStep().dependOn(&install_exe.step);

    //zcc.createStep(b, "cdb", targets.toOwnedSlice() catch @panic("OOM"));
}
