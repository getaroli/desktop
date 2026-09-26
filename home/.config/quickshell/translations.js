// Single translation table consumed by I18n.qml. Language payloads remain in
// their own files so translators can review one locale at a time.
.pragma library
.import "translations-es.js" as Es
.import "translations-pt-BR.js" as PtBR

var table = {
    "es": Es.es,
    "pt-BR": PtBR.ptBR
}
