// ControlPanel.qml, the control center, unfolded FROM the notch.
//
// Replaces swaync's "control center" (a GTK window of its own, hence
// impossible to fit in here). Left: vertical player. Center: sliders and
// toggles. Right: notification history, served by our own NotificationServer
// (see ShellState.qml).
import Quickshell
import Quickshell.Services.Notifications
import QtQuick
import QtQuick.Controls
import QtQuick.Layouts

Item {
    id: root

    // swallows clicks in the gap so they never close the panel
    MouseArea { anchors.fill: parent }

    // A short accent threshold ties the expanded control surface back to the
    // notch without turning the entire panel into a wallpaper-colored slab.
    Rectangle {
        anchors { top: parent.top; topMargin: 6; horizontalCenter: parent.horizontalCenter }
        width: 44
        height: 3
        radius: Appearance.radPill
        color: Colors.accent
        opacity: 0.72
    }

    RowLayout {
        anchors { fill: parent; leftMargin: 22; rightMargin: 22; topMargin: 20; bottomMargin: 18 }
        spacing: 20

        // ═════════════ left: vertical player ═════════════
        // It is the reference player, but INSIDE Super+D: sharing
        // the same surface and the same click-outside close as the rest of
        // the control center.
        MediaPanel {
            Layout.preferredWidth: 260
            Layout.fillHeight: true
        }

        Rectangle { Layout.fillHeight: true; width: 1; color: Qt.rgba(Colors.inkHi.r, Colors.inkHi.g, Colors.inkHi.b, 0.08) }

        // ═════════════════ center: controls ═════════════════
        ColumnLayout {
            Layout.preferredWidth: 410
            Layout.fillHeight: true
            spacing: 14

            // ---- sliders ----
            NotchSlider {
                Layout.fillWidth: true
                icon: ShellState.muted ? Icons.volMute : Icons.volHigh
                value: ShellState.muted ? 0 : ShellState.vol
                accent: Colors.accent
                onMoved: function (v) { ShellState.setVolume(v); }
                onIconClicked: ShellState.toggleMute()
            }
            // On a box with no internal panel (a tower) there is no
            // /sys/class/backlight, brightnessctl finds nothing and `bright`
            // stays -1 forever. The slider still showed: at zero, stuck and
            // with no effect. A knob commanding nothing lies about what you
            // can do, so the whole row goes and volume stays alone. Being a
            // ColumnLayout, an invisible item takes no room and the gap closes
            // on its own (same rule as Settings, which already hid its
            // brightness row with `shown`).
            NotchSlider {
                visible: ShellState.bright >= 0
                Layout.fillWidth: true
                icon: Icons.brightness
                value: Math.max(0, ShellState.bright)
                accent: Colors.warn
                onMoved: function (v) { ShellState.setBrightness(v); }
            }

            // ---- toggles ----
            GridLayout {
                Layout.fillWidth: true
                Layout.topMargin: 2
                columns: 5
                columnSpacing: 10

                // One target per tile. With a dedicated panel, any click
                // opens that panel; its switch already lives in the Net/
                // Bluetooth header. The rest of the toggles do act right here.
                component Toggle: Rectangle {
                    id: tg
                    property string icon: ""
                    property string label: ""
                    property bool on: false
                    property string panel: ""      // "" = no dedicated panel
                    property int idx: 0            // its turn in the cascade
                    signal activated()

                    Layout.fillWidth: true
                    Layout.preferredHeight: 58

                    // CASCADE. All six appeared at once, which is the same as
                    // saying none appeared: a block switching on has no
                    // direction and cannot be followed by eye. Staggered, the
                    // eye walks them left to right and the panel reads as
                    // something unfolding.
                    //
                    // NOTE: opacity stays untouched here. The fade is already
                    // done by the NotchLayer holding this panel, and two chained
                    // opacities multiply: not smoother, visibly late. The layer
                    // handles appearing; the cascade, only placing. One gesture,
                    // one owner.
                    property real ent: ShellState.mode === "control" ? 1 : 0
                    Behavior on ent {
                        SequentialAnimation {
                            PauseAnimation { duration: ShellState.mode === "control" ? tg.idx * Appearance.mStagger : 0 }
                            SpringAnimation {
                                spring: Appearance.sprPanel
                                damping: Appearance.dmpPanel
                                epsilon: Appearance.eppScale
                            }
                        }
                    }
                    scale: 0.82 + 0.18 * tg.ent
                    transform: Translate { y: (1 - tg.ent) * 12 }
                    radius: Appearance.radM
                    color: tg.on ? Qt.rgba(Colors.accent.r, Colors.accent.g, Colors.accent.b, 0.26)
                         : tgMa.containsMouse ? Qt.rgba(Colors.inkHi.r, Colors.inkHi.g, Colors.inkHi.b, 0.10)
                         : Qt.rgba(Colors.inkHi.r, Colors.inkHi.g, Colors.inkHi.b, 0.035)
                    border.width: 1
                    border.color: tg.on
                        ? Qt.rgba(Colors.accent.r, Colors.accent.g, Colors.accent.b, 0.46)
                        : tgMa.containsMouse
                            ? Qt.rgba(Colors.inkHi.r, Colors.inkHi.g, Colors.inkHi.b, 0.14)
                            : Qt.rgba(Colors.inkHi.r, Colors.inkHi.g, Colors.inkHi.b, 0.055)
                    Behavior on color { ColorAnimation { duration: Appearance.mQuick; easing.type: Easing.OutQuad } }
                    Behavior on border.color { ColorAnimation { duration: Appearance.mQuick; easing.type: Easing.OutQuad } }

                    ColumnLayout {
                        anchors.centerIn: parent
                        spacing: 1
                        Text {
                            Layout.alignment: Qt.AlignHCenter
                            text: tg.icon
                            color: tg.on ? Colors.accent : Colors.inkMid
                            font.family: Appearance.font; font.pixelSize: Appearance.fsL
                        }
                        Text {
                            Layout.alignment: Qt.AlignHCenter
                            text: tg.label
                            color: tg.on ? Colors.accent : Colors.inkMid
                            font.family: Appearance.fontUI; font.pixelSize: Appearance.fsCaption
                        }
                    }
                    MouseArea {
                        id: tgMa
                        anchors.fill: parent
                        hoverEnabled: true
                        cursorShape: Qt.PointingHandCursor
                        onClicked: {
                            if (tg.panel.length > 0) ShellState.togglePanel(tg.panel);
                            else tg.activated();
                        }
                    }
                }

                Toggle {
                    icon: ShellState.netIcon; label: I18n.tr("Network"); idx: 0
                    on: ShellState.online
                    panel: "network"
                }
                Toggle {
                    icon: ShellState.btIcon; label: I18n.tr("Bluetooth"); idx: 1
                    on: ShellState.btOn
                    panel: "bluetooth"
                }
                Toggle {
                    icon: ShellState.dnd ? "󰂛" : "󰂚"; label: I18n.tr("Do not disturb"); idx: 2
                    on: ShellState.dnd
                    onActivated: ShellState.dnd = !ShellState.dnd
                }
                Toggle {
                    icon: Icons.coffee; label: I18n.tr("Caffeine"); idx: 3
                    on: ShellState.caffeine
                    onActivated: ShellState.caffeine = !ShellState.caffeine
                }
                Toggle {
                    icon: Icons.moon; label: I18n.tr("Night light"); idx: 4
                    on: ShellState.nightLight
                    onActivated: ShellState.toggleNightLight()
                }
                Toggle {
                    icon: Icons.remote; label: I18n.tr("Remote"); idx: 5
                    on: ShellState.remoteMode
                    onActivated: ShellState.toggleRemoteMode()
                }
                Toggle {
                    icon: Icons.pokeball; label: I18n.tr("Pokémon"); idx: 6
                    on: ShellState.pokeTheme
                    onActivated: ShellState.togglePokeTheme()
                }
                Toggle {
                    icon: Icons.gamepad; label: I18n.tr("Game mode"); idx: 7
                    on: ShellState.gameMode
                    onActivated: ShellState.toggleGameMode()
                }
                Toggle {
                    icon: Icons.leaf; label: ShellState.batteryProfileLabel(ShellState.batteryProfile); idx: 8
                    visible: ShellState.batt >= 0 && ShellState.batteryProfileAvailable
                    on: ShellState.batteryProfile !== "balanced"
                    onActivated: ShellState.nextBatteryProfile()
                }
                // `on` always goes false on purpose: the other tiles switch
                // something on and color says whether it is on, but this one is
                // a door, not a switch. Painting it lit would promise a state
                // that does not exist.
                Toggle {
                    icon: Icons.calendar; label: I18n.tr("Calendar"); idx: 9
                    on: false
                    panel: "calendar"
                }
            }

            Item { Layout.fillHeight: true }

            // ---- time + system (inherited from the retired Sidebar) ----
            RowLayout {
                Layout.fillWidth: true
                Layout.topMargin: 2
                spacing: 16

                Text {
                    text: Icons.brightness
                    color: Colors.warn
                    font.family: Appearance.font; font.pixelSize: 14
                }
                Text {
                    Layout.fillWidth: true
                    text: ShellState.weather
                    color: Colors.inkMid; elide: Text.ElideRight
                    font.family: Appearance.fontUI; font.pixelSize: Appearance.fsXS
                }
                // Three glyphs with percentages were unreadable without knowing
                // by heart what each meant. Now this is a named door: keeps the
                // quick glance and opens the graphs.
                Rectangle {
                    id: performanceLink
                    Layout.preferredWidth: 222
                    Layout.preferredHeight: 42
                    radius: Appearance.radM
                    color: performanceMa.containsMouse
                        ? Qt.rgba(Colors.accent.r, Colors.accent.g, Colors.accent.b, 0.14)
                        : Qt.rgba(1, 1, 1, 0.05)
                    border.width: 1
                    border.color: performanceMa.containsMouse
                        ? Qt.rgba(Colors.accent.r, Colors.accent.g, Colors.accent.b, 0.26)
                        : Qt.rgba(1, 1, 1, 0.06)
                    Behavior on color { ColorAnimation { duration: Appearance.mQuick; easing.type: Easing.OutQuad } }

                    RowLayout {
                        anchors { fill: parent; leftMargin: 12; rightMargin: 10 }
                        spacing: 8
                        ColumnLayout {
                            Layout.fillWidth: true
                            spacing: 0
                            Text {
                                text: I18n.tr("Your computer")
                                color: performanceMa.containsMouse ? Colors.accent : Colors.inkMid
                                font.family: Appearance.fontUI; font.pixelSize: Appearance.fsCaption
                                font.weight: Font.DemiBold
                            }
                            Text {
                                text: I18n.tr("CPU {0}%  ·  RAM {1}%  ·  SSD {2}%",
                                              ShellState.cpu, ShellState.mem, ShellState.disk)
                                color: Colors.inkLo
                                font.family: Appearance.fontUI; font.pixelSize: Appearance.fsCaption
                                font.features: ({ "tnum": 1 })
                            }
                        }
                        Text {
                            text: "›"
                            color: performanceMa.containsMouse ? Colors.accent : Colors.inkLo
                            font.family: Appearance.fontUI; font.pixelSize: Appearance.fsTitle
                        }
                    }

                    MouseArea {
                        id: performanceMa
                        anchors.fill: parent
                        hoverEnabled: true
                        cursorShape: Qt.PointingHandCursor
                        onClicked: ShellState.togglePanel("system")
                    }
                }
            }
        }

        Rectangle { Layout.fillHeight: true; width: 1; color: Qt.rgba(Colors.inkHi.r, Colors.inkHi.g, Colors.inkHi.b, 0.08) }

        // ══════════════════ right: notifications ══════════════════
        ColumnLayout {
            Layout.fillWidth: true
            Layout.fillHeight: true
            spacing: 8

            RowLayout {
                Layout.fillWidth: true
                spacing: 8
                Text {
                    text: I18n.tr("Notifications")
                    color: Colors.inkHi
                    font.family: Appearance.fontUI; font.pixelSize: Appearance.fsS; font.weight: Font.DemiBold
                }
                Rectangle {
                    visible: ShellState.notifCount > 0
                    opacity: visible ? 1 : 0
                    Behavior on opacity { NumberAnimation { duration: Appearance.mIn; easing.type: Easing.OutCubic } }
                    implicitWidth: cnt.implicitWidth + 12
                    implicitHeight: 17
                    radius: Appearance.radPill
                    color: Qt.rgba(Colors.accent.r, Colors.accent.g, Colors.accent.b, 0.28)
                    Text {
                        id: cnt
                        anchors.centerIn: parent
                        text: ShellState.notifCount
                        color: Colors.accent
                        font.family: Appearance.fontUI; font.pixelSize: Appearance.fsCaption; font.weight: Font.Medium
                    }
                }
                Item { Layout.fillWidth: true }
                Text {
                    text: "󰒓"
                    color: gearMa.containsMouse ? Colors.accent : Colors.inkLo
                    font.family: Appearance.font; font.pixelSize: Appearance.fsL
                    Behavior on color { ColorAnimation { duration: Appearance.mQuick; easing.type: Easing.OutQuad } }
                    MouseArea {
                        id: gearMa
                        anchors.fill: parent; anchors.margins: -Appearance.hitPad
                        hoverEnabled: true; cursorShape: Qt.PointingHandCursor
                        onClicked: { ShellState.closePanel(); ShellState.settingsOpen = true; }
                    }
                }
                Text {
                    visible: ShellState.notifCount > 0
                    opacity: visible ? 1 : 0
                    Behavior on opacity { NumberAnimation { duration: Appearance.mIn; easing.type: Easing.OutCubic } }
                    text: I18n.tr("Clear")
                    color: clearMa.containsMouse ? Colors.accent : Colors.inkLo
                    font.family: Appearance.fontUI; font.pixelSize: Appearance.fsXS
                    Behavior on color { ColorAnimation { duration: Appearance.mQuick; easing.type: Easing.OutQuad } }
                    MouseArea {
                        id: clearMa
                        anchors.fill: parent; anchors.margins: -Appearance.hitPad
                        hoverEnabled: true; cursorShape: Qt.PointingHandCursor
                        onClicked: ShellState.clearNotifs()
                    }
                }
            }

            // ─────────── the mouse reads ONCE, and reads here ───────────
            // Each row had its MouseArea and the x another on top. Three bugs
            // came out, all felt right when closing one:
            //
            //   · the body goes in StyledText (whatever markup the app sends),
            //     and a Text with markup ACCEPTS hover events to detect links:
            //     it ate the row's. Measured: over the title the row lit and the
            //     x showed; two lines down, nothing. Most of each row's surface
            //     was dead zone.
            //   · the x appeared with `visible`, so it joined and left the
            //     layout: the text column widened 23 px the moment the mouse
            //     left it. The smallest of the three, but geometry moving alone
            //     right under the pointer.
            //   · and discarding one, the model is a NEW array (see
            //     ShellState.notifications): every delegate rebuilds with hover
            //     at zero. The hand never moved so no event arrives to correct
            //     it, and the list stays lying: x lit over a dark row, and the
            //     second click in the same spot closes nothing, landing in the
            //     gap.
            //
            // With a single MouseArea above the list one owner holds "what sits
            // under the pointer". And, this is what fixes the third, that owner
            // can look again when the list changes, without waiting for the hand
            // to move.
            Item {
                Layout.fillWidth: true
                Layout.fillHeight: true

                ListView {
                    id: notifList
                    anchors.fill: parent
                    clip: true
                    spacing: 7
                    model: ShellState.notifications
                    boundsBehavior: Flickable.StopAtBounds
                    ScrollBar.vertical: ScrollBar { policy: ScrollBar.AsNeeded; width: 4 }

                    property int hoverIdx: -1        // row under the pointer; -1 = none
                    property bool overClose: false   // ...and over its x

                    // (mx, my) in viewer coordinates, which is how the mouse reports them.
                    function track(mx, my) {
                        const y = my + notifList.contentY;
                        const i = notifList.indexAt(notifList.width / 2, y);
                        notifList.hoverIdx = i;
                        const fila = i >= 0 ? notifList.itemAtIndex(i) : null;
                        const equis = fila ? fila.closeItem : null;
                        if (!equis) {
                            notifList.overClose = false;
                            return;
                        }
                        // The target is never handwritten: the layout is asked
                        // where it left the x, plus the 6 px margin the
                        // MouseArea living there had. If the row spacing ever
                        // changes, the target moves alone.
                        const p = equis.mapToItem(fila, 0, 0);
                        notifList.overClose = mx >= p.x - 6 && mx <= p.x + equis.width + 6
                            && (y - fila.y) >= p.y - 6 && (y - fila.y) <= p.y + equis.height + 6;
                    }

                    // Looking again with no mouse move: closing one slides the
                    // ones below up and whoever sits under the pointer is
                    // already another row. forceLayout, or the list measures
                    // stale and the x lights on the wrong row.
                    function retrack() {
                        if (!listMa.containsMouse) {
                            notifList.hoverIdx = -1;
                            notifList.overClose = false;
                            return;
                        }
                        notifList.forceLayout();
                        notifList.track(listMa.mouseX, listMa.mouseY);
                    }

                    onCountChanged: Qt.callLater(notifList.retrack)
                    // With the wheel the rows move, not the hand.
                    onContentYChanged: if (listMa.containsMouse) notifList.track(listMa.mouseX, listMa.mouseY)

                    delegate: Rectangle {
                        id: nrow
                        required property var modelData
                        required property int index
                        readonly property bool hovered: notifList.hoverIdx === nrow.index
                        // All the row tells the list: where its x sits,
                        // so it knows who the click is for.
                        readonly property Item closeItem: equis

                        width: notifList.width
                        height: ncol.implicitHeight + 20
                        radius: Appearance.radS
                        color: nrow.hovered ? Qt.rgba(1, 1, 1, 0.08) : Qt.rgba(1, 1, 1, 0.045)
                        Behavior on color { ColorAnimation { duration: Appearance.mQuick; easing.type: Easing.OutQuad } }

                        RowLayout {
                            anchors { fill: parent; leftMargin: 12; rightMargin: 10; topMargin: 10; bottomMargin: 10 }
                            spacing: 11

                            Image {
                                Layout.preferredWidth: 22
                                Layout.preferredHeight: 22
                                Layout.alignment: Qt.AlignTop
                                source: nrow.modelData.image && nrow.modelData.image.length > 0
                                    ? nrow.modelData.image
                                    : Quickshell.iconPath(nrow.modelData.appIcon, "dialog-information")
                                sourceSize.width: 22; sourceSize.height: 22
                                fillMode: Image.PreserveAspectFit
                            }

                            ColumnLayout {
                                id: ncol
                                Layout.fillWidth: true
                                spacing: 2
                                Text {
                                    Layout.fillWidth: true
                                    text: (nrow.modelData.appName || "").toUpperCase()
                                    color: nrow.modelData.urgency === NotificationUrgency.Critical ? Colors.crit : Colors.inkMid
                                    elide: Text.ElideRight
                                    font.family: Appearance.fontUI; font.pixelSize: Appearance.fsCaption
                                    font.weight: Font.DemiBold
                                    font.letterSpacing: 0.5
                                }
                                Text {
                                    Layout.fillWidth: true
                                    text: nrow.modelData.summary || ""
                                    color: Colors.inkHi
                                    elide: Text.ElideRight
                                    font.family: Appearance.fontUI; font.pixelSize: Appearance.fsS; font.weight: Font.Medium
                                }
                                Text {
                                    Layout.fillWidth: true
                                    visible: text.length > 0
                                    text: nrow.modelData.body || ""
                                    color: Colors.inkMid
                                    wrapMode: Text.Wrap
                                    maximumLineCount: 3
                                    elide: Text.ElideRight
                                    textFormat: Text.StyledText
                                    font.family: Appearance.fontUI; font.pixelSize: Appearance.fsXS
                                }
                            }

                            Text {
                                id: equis
                                Layout.alignment: Qt.AlignTop
                                // Opacity, not `visible`: the x holds its slot lit
                                // or not, so the row measures the same with and
                                // without mouse. Reserving the seat costs nothing
                                // and removes the failure class where what you
                                // chase moves because you chase it.
                                opacity: nrow.hovered ? 1 : 0
                                text: "󰅖"
                                color: nrow.hovered && notifList.overClose ? Colors.crit : Colors.inkLo
                                font.family: Appearance.font; font.pixelSize: Appearance.fsS
                                Behavior on opacity { NumberAnimation { duration: Appearance.mQuick; easing.type: Easing.OutCubic } }
                                Behavior on color { ColorAnimation { duration: Appearance.mQuick; easing.type: Easing.OutQuad } }
                            }
                        }
                    }
                }

                MouseArea {
                    id: listMa
                    anchors.fill: parent
                    z: 1                     // above the rows: hover is its own
                    hoverEnabled: true
                    acceptedButtons: Qt.LeftButton
                    cursorShape: notifList.overClose ? Qt.PointingHandCursor : Qt.ArrowCursor
                    onPositionChanged: function (mouse) { notifList.track(mouse.x, mouse.y); }
                    onEntered: notifList.track(listMa.mouseX, listMa.mouseY)
                    onExited: { notifList.hoverIdx = -1; notifList.overClose = false; }
                    onClicked: {
                        const n = notifList.hoverIdx >= 0 ? notifList.model[notifList.hoverIdx] : null;
                        if (!n) return;
                        if (notifList.overClose) { ShellState.dismissNotif(n); return; }
                        const acts = n.actions;
                        if (acts && acts.length > 0) { acts[0].invoke(); ShellState.closePanel(); }
                    }
                }
            }

            Text {
                Layout.fillWidth: true
                Layout.fillHeight: true
                visible: ShellState.notifCount === 0
                opacity: visible ? 1 : 0
                Behavior on opacity { NumberAnimation { duration: Appearance.mIn; easing.type: Easing.OutCubic } }
                text: I18n.tr("No notifications")
                color: Colors.inkLo
                horizontalAlignment: Text.AlignHCenter
                verticalAlignment: Text.AlignVCenter
                font.family: Appearance.fontUI; font.pixelSize: Appearance.fsS
            }
        }
    }
}
