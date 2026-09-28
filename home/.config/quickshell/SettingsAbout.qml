// SettingsAbout.qml, system information. A single shell pass, not a
// poll: this does not change while you look at the window (except uptime,
// which refreshes when you re-enter the section).
//
// Each datum also carries its value as a `hint`, so if the text does not fit
// the row (the processor name, say) it shows whole in the footer strip on
// hover.
import Quickshell
import Quickshell.Io
import QtQuick
import QtQuick.Controls
import QtQuick.Effects
import QtQuick.Layouts

Flickable {
    id: root

    property string note: I18n.tr("One single pass when you open the section; no polling loop.")
    readonly property int matchCount: cSw.visibleRows + cHw.visibleRows + cAroli.visibleRows

    contentHeight: col.implicitHeight + 34
    clip: true
    boundsBehavior: Flickable.StopAtBounds
    ScrollBar.vertical: ScrollBar { policy: ScrollBar.AsNeeded; width: 5 }

    property var info: ({})

    onVisibleChanged: {
        if (!visible) return;
        contentY = 0;
        probe.running = true;
    }

    Process {
        id: probe
        command: ["bash", "-lc", `
            printf 'distro\\t%s\\n' "$(. /etc/os-release 2>/dev/null; echo "$PRETTY_NAME")"
            printf 'kernel\\t%s\\n' "$(uname -r)"
            printf 'wm\\t%s\\n'     "Hyprland $(hyprctl version 2>/dev/null | head -1 | grep -oP '\\d+\\.\\d+\\.\\d+' | head -1)"
            printf 'shell\\t%s\\n'  "Quickshell $(quickshell --version 2>/dev/null | grep -oP '\\d+\\.\\d+\\.\\d+' | head -1)"
            printf 'cpu\\t%s\\n'    "$(awk -F': ' '/model name/{print $2; exit}' /proc/cpuinfo)"
            printf 'ram\\t%s\\n'    "$(free -h --si | awk '/^Mem:/{print $3" ${I18n.tr("of")} "$2}')"
            printf 'disk\\t%s\\n'   "$(df -h --output=used,size,pcent / | tail -1 | awk '{print $1" ${I18n.tr("of")} "$2" ("$3")"}')"
            printf 'up\\t%s\\n'     "$(awk '{s=int($1); d=int(s/86400); h=int(s%86400/3600); m=int(s%3600/60);
                                          if(d>0) printf "%d d %d h %d min", d, h, m;
                                          else if(h>0) printf "%d h %d min", h, m;
                                          else printf "%d min", m}' /proc/uptime)"
            printf 'pkgs\\t%s\\n'   "$(pacman -Qq 2>/dev/null | wc -l) ${I18n.tr("packages")}"
            printf 'aroli\\t%s\\n'   "$("$HOME/.local/bin/aroli" version 2>/dev/null || echo "?")"
            printf 'rel\\t%s\\n'    "$(jq -r 'if .update_available then ((.remote // \"?\") + \" (!)\") else (.remote // \"?\") end' ~/.cache/aroli-desktop/update.json 2>/dev/null || echo "?")"
        `]
        stdout: SplitParser {
            onRead: function (line) {
                const p = line.split("\t");
                if (p.length < 2) return;
                const o = root.info;
                o[p[0]] = p[1];
                root.info = o;
                root.infoChanged();
            }
        }
    }

    readonly property var hw: [
        { k: "kernel", label: I18n.tr("Kernel") },
        { k: "cpu",    label: I18n.tr("Processor") },
        { k: "ram",    label: I18n.tr("Memory") },
        { k: "disk",   label: I18n.tr("Root disk") },
        { k: "up",     label: I18n.tr("Uptime") },
        { k: "pkgs",   label: I18n.tr("Installed packages") }
    ]

    readonly property var sw: [
        { k: "distro", label: I18n.tr("System") },
        { k: "wm",     label: I18n.tr("Compositor") },
        { k: "shell",  label: I18n.tr("Shell") }
    ]

    ColumnLayout {
        id: col
        width: root.width - 48
        x: 24
        y: 16
        spacing: 10

        // ─────────────────── header with logo ───────────────────
        RowLayout {
            Layout.fillWidth: true
            Layout.bottomMargin: 4
            spacing: 18
            visible: ShellState.settingsQuery.length === 0

            Text {
                visible: !Distro.omarchy
                text: Distro.glyph
                color: Colors.accent
                font.family: Appearance.font
                font.pixelSize: Appearance.fsXXL
            }
            Image {
                visible: Distro.omarchy
                source: Distro.markSource
                Layout.preferredWidth: 46
                Layout.preferredHeight: 46
                fillMode: Image.PreserveAspectFit
                // The mark is a grayscale SVG: tint it like the bar logo so
                // it follows the wallpaper palette.
                layer.enabled: true
                layer.effect: MultiEffect {
                    colorization: 1.0
                    colorizationColor: Colors.accent
                }
            }
            ColumnLayout {
                spacing: 2
                Text {
                    text: root.info["distro"] || Distro.name
                    color: Colors.inkHi
                    font.family: Appearance.fontUI
                    font.pixelSize: Appearance.fsTitle
                    font.weight: Font.DemiBold
                }
                Text {
                    text: root.info["wm"] || ""
                    color: Colors.inkMid
                    font.family: Appearance.fontUI
                    font.pixelSize: Appearance.fsS
                }
            }
        }

        SettingsControls.Card_ {
            id: cSw
            title: I18n.tr("SOFTWARE")
            Repeater {
                model: root.sw
                onItemAdded: cSw.recount()
                SettingsControls.Row_ {
                    required property var modelData
                    label: modelData.label
                    hint: root.info[modelData.k] || ""
                    SettingsControls.Val_ { text: root.info[modelData.k] || "…" }
                }
            }
        }

        SettingsControls.Card_ {
            id: cHw
            title: I18n.tr("HARDWARE")
            Repeater {
                model: root.hw
                onItemAdded: cHw.recount()
                SettingsControls.Row_ {
                    required property var modelData
                    label: modelData.label
                    hint: root.info[modelData.k] || ""
                    SettingsControls.Val_ { text: root.info[modelData.k] || "…" }
                }
            }
        }

        SettingsControls.Card_ {
            id: cAroli
            title: I18n.tr("THIS AROLI")

            SettingsControls.Row_ {
                label: I18n.tr("Aroli version")
                hint: root.info["aroli"] || ""
                SettingsControls.Val_ { text: root.info["aroli"] || "…" }
            }
            SettingsControls.Row_ {
                label: I18n.tr("Latest release seen")
                hint: I18n.tr("Written by the daily timer in ~/.cache/aroli-desktop/update.json. “(!)” means an update is available: run “aroli update --dry-run” in a terminal.")
                SettingsControls.Val_ { text: root.info["rel"] || "…" }
            }

            SettingsControls.Row_ {
                label: I18n.tr("Settings saved in")
                hint: I18n.tr("The JSON the Appearance section writes. Delete it and everything goes back to factory settings.")
                SettingsControls.Val_ { text: "~/.config/quickshell-rice.json" }
            }
            SettingsControls.Row_ {
                label: I18n.tr("Shell code")
                hint: I18n.tr("Bar, notch, panels, launcher, wallpaper picker and this very window.")
                SettingsControls.Val_ { text: "~/.config/quickshell/" }
            }
            SettingsControls.Row_ {
                label: I18n.tr("Motion language")
                hint: I18n.tr("The spec for durations, curves and shape that Hyprland and Quickshell follow.")
                SettingsControls.Val_ { text: "~/.config/motion-language.md" }
            }
        }

        SettingsControls.Note_ {
            Layout.topMargin: 10
            visible: ShellState.settingsQuery.length > 0 && !cSw.visible && !cHw.visible && !cAroli.visible
            text: I18n.tr("Nothing in About matches “{0}”.", ShellState.settingsQuery)
        }
    }
}
