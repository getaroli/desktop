// I18n.qml, the shell language layer.
//
// ENGLISH (EN-US) IS THE CODE. Strings stay written in English inside each
// .qml, wrapped in I18n.tr(...). In English tr() returns its argument as-is:
// no dictionary to maintain, no keys to invent, and no way for a text to come
// out untranslated and show up blank. Only Spanish and Brazilian Portuguese
// live in dictionaries (translations-es.js, translations-pt-BR.js), and
// whatever either one is missing falls back to English, which is a visible
// but harmless failure.
//
// WHY NOT qsTr() AND .qm: the Qt Linguist flow forces you to compile the .ts
// files and to restart the application to change language. Here the language
// switches from Settings and the whole interface repaints in place, because
// `lang` is a property and every `text: I18n.tr(...)` assignment is a binding
// that depends on it.
pragma Singleton

import Quickshell
import QtQuick
// The dictionaries go in imported .js files, not in another QML singleton:
// proven that a `readonly property var` holding all 400+ entries arrives as
// `undefined` when read from here. The full reason lives in the header of
// translations-es.js.
import "translations.js" as Translations

Singleton {
    id: root

    readonly property string lang: Config.language
    readonly property bool spanish: lang === "es"
    readonly property bool portugueseBrazil: lang === "pt-BR"

    // Each label goes in its own language on purpose: someone holding the
    // shell in a language they do not understand still has to be able to find
    // their own in the list.
    readonly property var languages: [
        { code: "pt-BR", label: "Português (Brasil)" },
        { code: "es", label: "Español" },
        { code: "en", label: "English" }
    ]

    function labelFor(code) {
        for (var i = 0; i < root.languages.length; i++)
            if (root.languages[i].code === code) return root.languages[i].label;
        return code;
    }

    function codeFor(label) {
        for (var i = 0; i < root.languages.length; i++)
            if (root.languages[i].label === label) return root.languages[i].code;
        return "pt-BR";
    }

    readonly property var labels: root.languages.map(function (l) { return l.label; })

    // tr("Shut down")                       -> "Apagar"
    // tr("{0} minutes left", 5)              -> "Quedan 5 minutos"
    //
    // The slots are {0}, {1}, {2} and are NEVER concatenated outside: a
    // sentence split into chunks ("Quedan " + n + " minutos") cannot be
    // translated, because in another language the pieces go in another
    // order.
    function tr(s, a0, a1, a2) {
        var dictionary = Translations.table[root.lang]
        var out = dictionary !== undefined && dictionary[s] !== undefined ? dictionary[s] : s;
        if (a0 !== undefined) out = out.split("{0}").join(a0);
        if (a1 !== undefined) out = out.split("{1}").join(a1);
        if (a2 !== undefined) out = out.split("{2}").join(a2);
        return out;
    }
}
