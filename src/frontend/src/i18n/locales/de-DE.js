/**
 * Deutsches Sprachpaket
 * @author Chen Fengqing
 * @date 2025-10-05
 */

export default {
  // Allgemein
  common: {
    confirm: "Bestätigen",
    cancel: "Abbrechen",
    save: "Speichern",
    delete: "Löschen",
    edit: "Bearbeiten",
    add: "Hinzufügen",
    search: "Suchen",
    close: "Schließen",
    ok: "OK",
    yes: "Ja",
    no: "Nein",
    loading: "Laden...",
    success: "Erfolg",
    error: "Fehler",
    warning: "Warnung",
    info: "Information",
    copy: "Kopieren",
    paste: "Einfügen",
    cut: "Ausschneiden",
    refresh: "Aktualisieren",
    reset: "Zurücksetzen",
    clear: "Löschen",
    back: "Zurück",
    next: "Weiter",
    previous: "Vorherige",
    finish: "Fertig",
    submit: "Senden",
    import: "Importieren",
    export: "Exportieren",
    settings: "Einstellungen",
    help: "Hilfe",
    about: "Über",
    more: "Mehr",
  },

  // Anmeldeseite
  login: {
    title: "Tresor-Anmeldung",
    password: "Passwort",
    passwordPlaceholder: "Tresor-Passwort eingeben",
    login: "Anmelden",
    createVault: "Tresor erstellen",
    openVault: "Anderen Tresor öffnen",
    invalidPassword: "Falsches Passwort",
    loginSuccess: "Anmeldung erfolgreich",
    loginFailed: "Anmeldung fehlgeschlagen",
  },

  // Hauptoberfläche
  main: {
    title: "wepass",
    searchPlaceholder: "Konten suchen...",
    noData: "Keine Daten",
    noSearchResults: "Keine passenden Ergebnisse gefunden",
  },

  // Einstellungsdialog
  settings: {
    title: "Einstellungen",
    general: "Allgemein",
    log: "Protokoll",
    lock: "Sperren",
    theme: "Design",
    language: "Sprache",
    lightTheme: "Hell",
    darkTheme: "Dunkel",
    selectTheme: "Design auswählen",
    selectLanguage: "Sprache auswählen",
    logSettings: "Protokolleinstellungen",
    enableInfoLog: "Info-Protokoll aktivieren",
    enableDebugLog: "Debug-Protokoll aktivieren",
    lockSettings: "Sperreinstellungen",
    enableAutoLock: "Automatische Sperre aktivieren",
    enableTimerLock: "Timer-Sperre aktivieren",
    enableMinimizeLock: "Sperre beim Minimieren aktivieren",
    lockTimeMinutes: "Sperrzeit (Minuten)",
    enableSystemLock: "System-Sperre aktivieren",
    systemLockMinutes: "System-Sperrzeit (Minuten)",
    settingsSaved: "Einstellungen gespeichert",
    loadSettingsFailed: "Laden der Einstellungen fehlgeschlagen",
    saveSettingsFailed: "Speichern der Einstellungen fehlgeschlagen",
  },

  // Kontoverwaltung
  account: {
    title: "Titel",
    username: "Benutzername",
    password: "Passwort",
    url: "URL",
    type: "Typ",
    notes: "Notizen",
    icon: "Symbol",
    favorite: "Favorit",
    useCount: "Verwendungsanzahl",
    lastUsed: "Zuletzt verwendet",
    created: "Erstellt",
    updated: "Aktualisiert",
    inputMethod: "Eingabemethode",
    addAccount: "Konto hinzufügen",
    editAccount: "Konto bearbeiten",
    deleteAccount: "Konto löschen",
    copyUsername: "Benutzername kopieren",
    copyPassword: "Passwort kopieren",
    openUrl: "URL öffnen",
    generatePassword: "Passwort generieren",
    showPassword: "Passwort anzeigen",
    hidePassword: "Passwort verbergen",
    accountSaved: "Konto gespeichert",
    accountDeleted: "Konto gelöscht",
    confirmDelete: "Sind Sie sicher, dass Sie dieses Konto löschen möchten?",
    deleteConfirmTitle: "Löschbestätigung",
    copiedToClipboard: "In die Zwischenablage kopiert",
    copyFailed: "Kopieren fehlgeschlagen",
  },

  // Gruppen und Typen
  group: {
    allGroups: "Alle Gruppen",
    addGroup: "Gruppe hinzufügen",
    editGroup: "Gruppe bearbeiten",
    deleteGroup: "Gruppe löschen",
    groupName: "Gruppenname",
    groupIcon: "Gruppensymbol",
    rename: "Umbenennen", // TODO: Translate
    createGroup: "Gruppe erstellen", // TODO: Translate
    moveLeft: "Nach links verschieben", // TODO: Translate
    moveRight: "Nach rechts verschieben", // TODO: Translate
    newGroupName: "Neuen Gruppennamen eingeben", // TODO: Translate
    renameGroup: "Gruppe umbenennen", // TODO: Translate
    newGroupPlaceholder: "Gruppenname", // TODO: Translate
    groupNameCannotBeEmpty: "Gruppenname darf nicht leer sein", // TODO: Translate
    groupNameNotChanged: "Gruppenname nicht geändert", // TODO: Translate
    defaultGroupCannotBeDeleted: "Standardgruppe kann nicht gelöscht werden", // TODO: Translate
    confirmDeleteGroup:
      'Sind Sie sicher, dass Sie die Gruppe "{groupName}" löschen möchten?\nAlle Tags und Konten in dieser Gruppe werden ebenfalls gelöscht, diese Aktion ist unumkehrbar.', // TODO: Translate
    deleteGroupTitle: "Gruppe löschen", // TODO: Translate
    confirmDelete: "Löschen bestätigen", // TODO: Translate
    defaultGroupName: "Standard", // TODO: Translate
  },

  type: {
    allTypes: "Alle Typen",
    addType: "Typ hinzufügen",
    editType: "Typ bearbeiten",
    deleteType: "Typ löschen",
    typeName: "Typname",
    typeIcon: "Typsymbol",
  },

  // Hilfedialog
  help: {
    title: "Hilfe",
    quickStart: "Schnellstart",
    welcome: "Willkommen bei WePassword",
    welcomeTitle: "Willkommen bei WePassword",
    description:
      "WePassword ist ein sicheres Passwort-Management-Tool, das Ihnen hilft, Passwörter sicher zu speichern und zu verwalten.",
    welcomeDescription:
      "WePassword ist ein sicheres Passwort-Management-Tool, das Ihnen hilft, Passwörter sicher zu speichern und zu verwalten.",
    basicOperations: "Grundlegende Operationen",
    createVault:
      "Tresor erstellen: Bei der ersten Verwendung führt Sie das System durch die Erstellung eines neuen Tresors",
    addPassword:
      'Passwort hinzufügen: Klicken Sie auf die Schaltfläche "Hinzufügen", um Website-, Benutzername- und Passwort-Informationen einzugeben',
    searchPassword:
      "Passwort suchen: Verwenden Sie das Suchfeld, um schnell bestimmte Passwort-Einträge zu finden",
    editPassword:
      "Passwort bearbeiten: Doppelklicken Sie auf einen Passwort-Eintrag oder klicken Sie auf die Bearbeiten-Schaltfläche zum Ändern",
    deletePassword:
      "Passwort löschen: Wählen Sie einen Passwort-Eintrag aus und klicken Sie auf die Löschen-Schaltfläche",
  },

  // Über-Dialog
  about: {
    title: "Über",
    appName: "Passwort-Manager",
    description:
      "Plattformübergreifendes Passwort-Management-Tool basierend auf Wails + Go + Vue.js",
    version: "Version",
    buildDate: "Build-Datum",
    author: "Autor",
    support: "Support",
    license: "Lizenz",
    github: "GitHub",
    gitee: "Gitee",
  },

  // Fehlermeldungen
  error: {
    networkError: "Netzwerkfehler",
    serverError: "Serverfehler",
    unknownError: "Unbekannter Fehler",
    operationFailed: "Operation fehlgeschlagen",
    dataLoadFailed: "Laden der Daten fehlgeschlagen",
    saveFailed: "Speichern fehlgeschlagen",
    deleteFailed: "Löschen fehlgeschlagen",
    copyFailed: "Kopieren fehlgeschlagen",
    invalidInput: "Ungültige Eingabe",
    requiredField: "Dieses Feld ist erforderlich",
  },

  // Erfolgsmeldungen
  success: {
    operationSuccess: "Operation erfolgreich",
    dataSaved: "Daten gespeichert",
    dataDeleted: "Daten gelöscht",
    copied: "Kopiert",
    imported: "Import erfolgreich",
    exported: "Export erfolgreich",
    dataRefreshed: "Daten aktualisiert",
    tabRenamed: "Tab erfolgreich umbenannt",
    tabDeleted: "Tab erfolgreich gelöscht",
    tabMoved: "Tab erfolgreich verschoben",
    tabCreated: "Tab erfolgreich erstellt",
    groupCreated: "Gruppe erfolgreich erstellt",
    groupRenamed: "Gruppe erfolgreich umbenannt",
    groupDeleted: "Gruppe erfolgreich gelöscht",
    groupMoved: "Gruppe erfolgreich verschoben",
    accountSaved: "Konto erfolgreich gespeichert",
    accountDeleted: "Konto erfolgreich gelöscht",
    passwordGenerated: "Passwort erfolgreich generiert",
    settingsSaved: "Einstellungen erfolgreich gespeichert",
    vaultExported: "Tresor erfolgreich exportiert",
    vaultImported: "Tresor erfolgreich importiert",
    usernameCopied:
      "Benutzername in Zwischenablage kopiert (wird nach 10 Sekunden automatisch gelöscht)", // TODO: Translate
    passwordCopied: "Passwort in Zwischenablage kopiert", // TODO: Translate
    usernameAndPasswordCopied:
      "Benutzername und Passwort in Zwischenablage kopiert (wird nach 10 Sekunden automatisch gelöscht)", // TODO: Translate
    urlCopied: "URL in Zwischenablage kopiert", // TODO: Translate
    titleCopied: "Titel in Zwischenablage kopiert", // TODO: Translate
    notesCopied:
      "Notizen in Zwischenablage kopiert (wird nach 10 Sekunden automatisch gelöscht)", // TODO: Translate
  },

  // Warnmeldungen
  warning: {
    noGroupData: "Keine Gruppendaten, bitte erstellen Sie zuerst eine Gruppe",
    selectGroupFirst: "Bitte wählen Sie zuerst eine Gruppe aus",
    selectTabFirst: "Bitte wählen Sie zuerst einen Tab aus",
    noAccountSelected: "Bitte wählen Sie zuerst ein Konto aus",
    confirmOperation: "Bitte bestätigen Sie diese Operation",
  },

  // Statusinformationen
  status: {
    loading: "Laden...",
    saving: "Speichern...",
    deleting: "Löschen...",
    processing: "Verarbeiten...",
    connecting: "Verbinden...",
    searching: "Suchen...",
    exporting: "Exportieren...",
    importing: "Importieren...",
    generating: "Generieren...",
  },

  // Mehr-Menü
  moreMenu: {
    more: "Mehr",
    selectNewVault: "Neuen Tresor auswählen",
    openVaultDirectory: "Tresor-Verzeichnis öffnen",
    generatePassword: "Passwort generieren",
    setPasswordRules: "Passwort-Generierungsregeln festlegen",
    changeLoginPassword: "Anmeldepasswort ändern",
    exportVault: "Tresor exportieren",
    importVault: "Tresor importieren",
    changeLog: "Änderungsprotokoll",
    settings: "Einstellungen",
    help: "Hilfe",
    about: "Über",
    lockVault: "Tresor sperren",
    logout: "Abmelden",
    oldLoginPasswordLabel: "Altes Anmeldepasswort",
    oldLoginPasswordPlaceholder: "Aktuelles Anmeldepasswort eingeben",
    newLoginPasswordLabel: "Neues Anmeldepasswort",
    newLoginPasswordPlaceholder: "Neues Anmeldepasswort eingeben",
    confirmNewPasswordLabel: "Neues Passwort bestätigen",
    confirmNewPasswordPlaceholder: "Neues Anmeldepasswort erneut eingeben",
    selectVaultFile: "Tresor-Datei auswählen",
    selectVaultFilePrompt: "Bitte wählen Sie die zu öffnende Tresor-Datei aus:",
    selectVaultFilePlaceholder: "Bitte Tresor-Dateipfad auswählen",
    browse: "Durchsuchen",
    supportedFileFormats: "Unterstützte Dateiformate: .db, .vault",
    openFileFailed: "Datei öffnen fehlgeschlagen",
    passwordCopied: "Passwort in Zwischenablage kopiert",
    copyFailedManual:
      "Kopieren fehlgeschlagen, bitte Passwort manuell kopieren",
    passwordChangeInProgress: "Passwort-Änderung läuft, bitte warten...",
    formRefNotFound: "Formular-Referenz nicht gefunden",
    oldPasswordIncorrect: "Altes Passwort falsch, bitte erneut eingeben",
    changePasswordConfirm:
      "Das Ändern des Anmeldepassworts verschlüsselt alle Kontodaten neu. Dieser Vorgang ist nicht rückgängig zu machen, stellen Sie sicher, dass Sie Ihren aktuellen Tresor gesichert haben. Fortfahren?",
    passwordChangeSuccess: "Anmeldepasswort erfolgreich geändert!",
    passwordChangeFailed: "Anmeldepasswort ändern fehlgeschlagen",
    importSuccessDataUpdated: "Import erfolgreich, Daten aktualisiert",
    importSuccessDataRefreshFailed:
      "Import erfolgreich, aber Datenaktualisierung fehlgeschlagen, bitte Seite manuell aktualisieren",
    vaultDirectoryOpened: "Tresor-Verzeichnis geöffnet",
    openDirectoryFailed: "Verzeichnis öffnen fehlgeschlagen",
    lockVaultConfirm:
      "Sind Sie sicher, dass Sie den Tresor sperren möchten? Nach dem Sperren müssen Sie das Passwort erneut eingeben, um darauf zuzugreifen.",
    vaultLocked: "Tresor gesperrt",
    lockVaultFailed: "Tresor sperren fehlgeschlagen",
    logoutConfirm:
      "Sind Sie sicher, dass Sie sich abmelden möchten? Nach der Abmeldung müssen Sie das Passwort erneut eingeben, um auf den Tresor zuzugreifen.",
    logoutSuccess: "Abmeldung erfolgreich",
    oldLoginPasswordRequired: "Bitte altes Passwort eingeben",
    newLoginPasswordRequired: "Bitte neues Passwort eingeben",
    newLoginPasswordMinLength: "Passwort muss mindestens 8 Zeichen lang sein",
    newLoginPasswordStrength:
      "Passwort muss Groß-, Kleinbuchstaben, Zahlen und Sonderzeichen enthalten",
    confirmNewLoginPasswordRequired: "Bitte neues Passwort bestätigen",
    passwordsDoNotMatch: "Passwörter stimmen nicht überein",
    selectNewVaultConfirm:
      "Die Auswahl eines neuen Tresors führt Sie zum Anmeldebildschirm zurück, Sie müssen die Tresor-Datei erneut auswählen und das Passwort eingeben. Fortfahren?",
    selectNewVaultSuccess:
      "Zum Anmeldebildschirm zurückgekehrt, bitte neuen Tresor auswählen",
    pleaseSelectFile: "Bitte zuerst Datei auswählen",
    openingFile: "Datei wird geöffnet: ",
    open: "Öffnen",
    continue: "Fortfahren",
    verifying: "Überprüfung...",
    changingPassword: "Änderung...",
  },

  // Kontextmenü
  contextMenu: {
    inputUsernameAndPassword: "Benutzernamen und Passwort eingeben", // TODO: Translate
    openUrl: "URL öffnen",
    duplicate: "Kopie erstellen",
    view: "Anzeigen",
    edit: "Bearbeiten",
    changeGroup: "Gruppe ändern",
    copyUsername: "Benutzername kopieren",
    copyPassword: "Passwort kopieren",
    copyUsernameAndPassword: "Benutzername und Passwort kopieren",
    copyUrl: "URL kopieren",
    copyTitle: "Titel kopieren",
    copyNotes: "Notizen kopieren",
    delete: "Löschen",
  },

  // Export-Funktion
  export: {
    title: "Tresor exportieren",
    steps: {
      verifyPassword: "Passwort verifizieren",
      selectAccounts: "Konten auswählen",
      setBackup: "Backup einrichten",
      exportComplete: "Export abgeschlossen",
    },
    verifyPasswordTitle: "Anmeldepasswort verifizieren",
    verifyPasswordDesc:
      "Bitte geben Sie das aktuelle Anmeldepasswort des Tresors ein, um den Export fortzusetzen.",
    loginPassword: "Anmeldepasswort",
    loginPasswordPlaceholder: "Bitte Anmeldepasswort eingeben",
    selectAccountsTitle: "Zu exportierende Konten auswählen",
    exportAll: "Alle exportieren",
    exportByGroup: "Nach Gruppe exportieren",
    exportByType: "Nach Typ exportieren", // TODO: Translate
    exportSelected: "Ausgewählte exportieren",
    selectGroups: "Gruppen auswählen", // TODO: Translate
    selectAll: "Alle auswählen", // TODO: Translate
    clearAll: "Auswahl aufheben", // TODO: Translate
    groupSelectionSummary:
      "Ausgewählte {count} Gruppen, erwartet werden {accountCount} Konten", // TODO: Translate
    selectTypes: "Typen auswählen", // TODO: Translate
    typeSelectionSummary:
      "Ausgewählte {count} Typen, erwartet werden {accountCount} Konten", // TODO: Translate
    selectAccounts: "Konten auswählen", // TODO: Translate
    loadingAccounts: "Konten werden geladen...", // TODO: Translate
    noAccounts: "Keine Kontodaten", // TODO: Translate
    setBackupTitle: "Backup-Passwort und Exportpfad festlegen", // TODO: Translate
    backupPassword: "Backup-Passwort", // TODO: Translate
    backupPasswordPlaceholder: "Bitte Backup-Passwort eingeben", // TODO: Translate
    generate: "Generieren", // TODO: Translate
    backupPasswordTip:
      "Das Backup-Passwort wird zum Verschlüsseln der exportierten Daten verwendet, bitte bewahren Sie es sicher auf", // TODO: Translate
    exportPath: "Exportpfad", // TODO: Translate
    exportPathPlaceholder: "Bitte Exportpfad auswählen", // TODO: Translate
    browse: "Durchsuchen", // TODO: Translate
    exporting: "Tresor wird exportiert...", // TODO: Translate
    exportSuccessTitle: "Export erfolgreich", // TODO: Translate
    exportSuccessSubTitle: "Tresor wurde erfolgreich nach: {path} exportiert", // TODO: Translate
    exportFailedTitle: "Export fehlgeschlagen", // TODO: Translate
    startExport: "Export starten", // TODO: Translate
    openFolder: "Ordner öffnen", // TODO: Translate
    loginPasswordRequired: "Bitte Anmeldepasswort eingeben", // TODO: Translate
    backupPasswordRequired: "Bitte Backup-Passwort eingeben", // TODO: Translate
    backupPasswordMinLength:
      "Backup-Passwort muss mindestens 6 Zeichen lang sein", // TODO: Translate
    exportPathRequired: "Bitte Exportpfad auswählen", // TODO: Translate
  },

  // Import-Funktion
  import: {
    title: "Tresor importieren",
    selectFile: "Datei auswählen",
    selectFileDesc: "Bitte wählen Sie die zu importierende Tresor-Datei aus",
    fileFormat: "Dateiformat",
    importProgress: "Import-Fortschritt",
    importComplete: "Import abgeschlossen",
  },

  // Passwort-Generator
  passwordGenerator: {
    title: "Passwort generieren",
    selectRule: "Passwort-Regel auswählen",
    selectRulePlaceholder: "Bitte Passwort-Regel auswählen",
    generalRule: "Allgemeine Passwort-Regel",
    customRule: "Benutzerdefinierte Passwort-Regel",
    includeUppercase: "Großbuchstaben einschließen",
    includeLowercase: "Kleinbuchstaben einschließen",
    includeNumbers: "Zahlen einschließen",
    includeSpecialChars: "Sonderzeichen einschließen",
    passwordLength: "Passwort-Länge",
    generatePassword: "Passwort generieren",
    copyPassword: "Passwort kopieren",
    regenerate: "Neu generieren",
    customSpecialChars: "Benutzerdefinierte Sonderzeichen",
    defaultSpecialChars: "Leer lassen für Standard-Sonderzeichen",
    usePassword: "Dieses Passwort verwenden",
    generatedPassword: "Generiertes Passwort",
    clickToGenerate: "Klicken Sie auf Passwort generieren",
    selectCharType: "Bitte wählen Sie mindestens einen Zeichentyp aus",
    enterPattern: "Bitte geben Sie ein Passwort-Muster ein",
    generateFirst: "Bitte generieren Sie zuerst ein Passwort",
  },

  // Passwort-Regel-Einstellungen
  passwordRules: {
    title: "Passwort-Regel-Einstellungen",
    ruleName: "Regelname",
    ruleNamePlaceholder: "Bitte Regelname eingeben",
    saveRule: "Regel speichern",
    deleteRule: "Regel löschen",
    editRule: "Regel bearbeiten",
  },

  // Test-bezogen
  test: {
    passwordDetailTest: "Passwort-Detail-Komponenten-Test",
  },

  // Import Vault Feature
  importVault: {
    title: "Passwort-Tresor importieren",
    step1: "Datei auswählen",
    step2: "Passwort verifizieren",
    step3: "Import abgeschlossen",
    step1Title: "Import-Datei auswählen",
    step1Description:
      "Bitte wählen Sie die zu importierende Passwort-Tresor-Backup-Datei (ZIP-Format) aus.",
    importFile: "Import-Datei",
    selectImportFilePlaceholder: "Bitte Import-Datei auswählen",
    browse: "Durchsuchen",
    step2Title: "Entpackungspasswort eingeben",
    step2Description:
      "Bitte geben Sie das Entpackungspasswort (Backup-Passwort) für die Datei ein.",
    backupPassword: "Entpackungspasswort",
    enterBackupPasswordPlaceholder: "Bitte Entpackungspasswort eingeben",
    importingVault: "Passwort-Tresor wird importiert...",
    importComplete: "Import abgeschlossen",
    importSuccess: "Import erfolgreich",
    vaultImportedSuccessfully: "Passwort-Tresor erfolgreich importiert",
    importReport: "Import-Bericht",
    totalAccounts: "Gesamtanzahl Konten",
    successfullyImported: "Erfolgreich importiert",
    skippedAccounts: "Übersprungene Konten",
    errorAccounts: "Fehlerhafte Konten",
    totalGroups: "Gesamtanzahl Gruppen",
    importedGroups: "Importierte Gruppen",
    totalTypes: "Gesamtanzahl Typen",
    importedTypes: "Importierte Typen",
    skippedAccountDetails: "Details übersprungener Konten",
    accountTitle: "Konto-Titel",
    accountName: "Konto-Name",
    accountId: "Konto-ID",
    importFailed: "Import fehlgeschlagen",
    importFailedStatus: "Import fehlgeschlagen",
    unknownError: "Unbekannter Fehler beim Import aufgetreten",
    close: "Schließen",
    cancel: "Abbrechen",
    previousStep: "Vorheriger Schritt",
    startImport: "Import starten",
    nextStep: "Nächster Schritt",
    refreshData: "Daten aktualisieren",
    selectImportFileMessage: "Bitte Import-Datei auswählen",
    enterBackupPasswordMessage: "Bitte Entpackungspasswort eingeben",
    selectFileSuccess: "Import-Datei erfolgreich ausgewählt",
    selectFileFailed: "Auswahl der Import-Datei fehlgeschlagen",
    preparingToImport: "Import wird vorbereitet...",
    validatingFileAndPassword: "Datei und Passwort werden validiert...",
    vaultImportSuccess: "Passwort-Tresor erfolgreich importiert",
    dataRefreshed: "Daten aktualisiert",
    importInProgressWarning: "Import läuft, bitte warten...",
  },

  // Passwort-Regel-Einstellungen
  passwordRuleSettings: {
    title: "Passwort-Generierungsregel-Einstellungen",
    savedRules: "Gespeicherte Passwort-Regeln",
    newRule: "Neue Regel",
    ruleName: "Regelname",
    description: "Beschreibung",
    operation: "Operation",
    edit: "Bearbeiten",
    delete: "Löschen",
    generalRule: "Allgemeine Passwort-Regel",
    includeUppercase: "Großbuchstaben einschließen",
    includeLowercase: "Kleinbuchstaben einschließen",
    includeNumbers: "Zahlen einschließen",
    includeSpecialChars: "Sonderzeichen einschließen",
    passwordLength: "Passwort-Länge",
    customSpecialChars: "Benutzerdefinierte Sonderzeichen",
    customSpecialCharsPlaceholder:
      "Leer lassen für Standard-Sonderzeichen: !@#$%^&*()_+-=[]{}|;:,.<>?",
    customRulesDescription: "Beschreibung benutzerdefinierter Passwort-Regeln",
    close: "Schließen",
    saveSettings: "Einstellungen speichern",
    editRule: "Regel bearbeiten",
    newRuleDialog: "Neue Regel",
    ruleNameRequired: "Bitte Regelname eingeben",
    ruleDescriptionPlaceholder: "Bitte Regelbeschreibung eingeben",
    ruleType: "Regeltyp",
    general: "Allgemeine Regel",
    custom: "Benutzerdefinierte Regel",
    passwordPattern: "Passwort-Muster",
    passwordPatternPlaceholder: "Beispiel: Aaa111",
    save: "Speichern",
    cancel: "Abbrechen",
    confirmDeleteRule:
      'Sind Sie sicher, dass Sie die Regel "{ruleName}" löschen möchten?',
    deleteConfirmation: "Löschbestätigung",
    confirmDelete: "Bestätigen",
    ruleDeletedSuccess: "Regel erfolgreich gelöscht",
    ruleUpdateSuccess: "Regel erfolgreich aktualisiert",
    ruleCreateSuccess: "Regel erfolgreich erstellt",
    enterRuleName: "Bitte Regelname eingeben",
    settingsSaved: "Einstellungen erfolgreich gespeichert",
  },

  // Suchergebnisse
  searchResults: {
    groups: "Gruppen",
    accounts: "Konten",
    noUrl: "Keine URL",
    belongsToGroup: "Gehört zur Gruppe",
    unknownGroup: "Unbekannte Gruppe",
    noSearchResults: "Keine Suchergebnisse gefunden",
    tryOtherKeywords: "Versuchen Sie andere Suchbegriffe",
  },

  // Statusleiste
  statusBar: {
    authInfo: "Autorisierungsinfo: Pro-Version aktiviert",
    usageCount: "Nutzung: {count} Mal",
    usageDays: "Nutzungstage: {days} Tage",
  },

  // Tab-Kontextmenü
  tabContextMenu: {
    rename: "Umbenennen",
    deleteTab: "Tab löschen",
    newTab: "Neuer Tab",
    moveUp: "Nach oben",
    moveDown: "Nach unten",
    promptNewName: "Bitte neuen Tab-Namen eingeben",
    renameTab: "Tab umbenennen",
    confirm: "Bestätigen",
    cancel: "Abbrechen",
    tabName: "Tab-Name",
    tabNameCannotBeEmpty: "Tab-Name darf nicht leer sein",
    tabNameNotChanged: "Tab-Name nicht geändert",
    userCanceledRename: "Benutzer hat Umbenennung abgebrochen",
    confirmDeleteTab:
      'Sind Sie sicher, dass Sie Tab "{tabName}" löschen möchten?\nNach dem Löschen werden auch alle Konten unter diesem Tab gelöscht, dieser Vorgang ist nicht rückgängig zu machen.',
    deleteConfirmation: "Tab löschen",
    confirmDelete: "Löschen bestätigen",
    userCanceledDelete: "Benutzer hat Löschvorgang abgebrochen",
    promptNewTabName: "Bitte neuen Tab-Namen eingeben",
    userCanceledNewTab: "Benutzer hat neuen Tab abgebrochen",
  },

  // Tabs-Seitenleiste
  tabsSidebar: {
    emptyTabs: "Keine Tabs",
    createTabHint:
      "Klicken Sie auf die Schaltfläche unten, um einen Tab zu erstellen",
    newTab: "Neuer Tab",
  },

  // Test-Dialog
  testDialog: {
    title: "Test-Dialog",
    description1:
      "Dies ist ein Test-Dialog zur Überprüfung, ob die Popup-Funktion ordnungsgemäß funktioniert.",
    description2:
      "Wenn Sie diesen Dialog sehen können, bedeutet das, dass die Popup-Funktion ordnungsgemäß funktioniert.",
    close: "Schließen",
  },

  // Titelleiste
  titleBar: {
    defaultAppTitle: "Passwort-Manager",
  },

  // Sperr-Event-Service
  lockEventService: {
    logPrefix: "Sperr-Event-Service",
    frontendStateUpdated: "Frontend-Status aktualisiert",
    redirectToLogin: "Weiterleitung zur Anmeldeseite",
    vaultAutoLocked:
      "Passwort-Tresor wurde automatisch gesperrt, bitte erneut anmelden",
    sensitiveDataCleared: "Sensible Daten gelöscht",
    clearSensitiveDataFailed: "Löschen sensibler Daten fehlgeschlagen",
    windowEventListenersSet: "Fenster-Event-Listener wurden gesetzt",
    windowLostFocus: "Fenster hat Fokus verloren",
    windowGainedFocus: "Fenster hat Fokus erhalten",
    backendNotifiedMinimize: "Backend über Fenster-Minimierung benachrichtigt",
    notifyMinimizeFailed:
      "Benachrichtigung über Fenster-Minimierung fehlgeschlagen",
    backendNotifiedFocus: "Backend über Fenster-Fokus benachrichtigt",
    notifyFocusFailed: "Benachrichtigung über Fenster-Fokus fehlgeschlagen",
    manualLockCheckTriggered: "Manuelle Sperrprüfung ausgelöst",
  },

  // Konto-Utilities
  accountUtils: {
    defaultAccountTitle: "Konto",
    startCopyingPassword: "Beginne Passwort zu kopieren, Konto-ID: {accountId}",
    passwordCopiedSuccess:
      "Passwort in Zwischenablage kopiert (automatische Löschung nach 10 Sekunden)",
    passwordCopySuccess: "Passwort erfolgreich kopiert, Konto-ID: {accountId}",
    passwordCopyFailed: "Passwort kopieren fehlgeschlagen",
    startGettingPassword: "Beginne Passwort zu holen, Konto-ID: {accountId}",
    passwordGetSuccess: "Passwort erfolgreich geholt, Konto-ID: {accountId}",
    passwordGetFailed: "Passwort holen fehlgeschlagen",
    startGettingAccountDetail:
      "Beginne Konto-Details zu holen, Konto-ID: {accountId}",
    accountDetailGetSuccess:
      "Konto-Details erfolgreich geholt, Konto-ID: {accountId}",
    accountDetailGetFailed: "Konto-Details holen fehlgeschlagen",
  },
};
