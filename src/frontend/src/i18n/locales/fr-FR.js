/**
 * Pack de langue française
 * @author Chen Fengqing
 * @date 2025-10-05
 */

export default {
  // Commun
  common: {
    confirm: "Confirmer",
    cancel: "Annuler",
    save: "Enregistrer",
    delete: "Supprimer",
    edit: "Modifier",
    add: "Ajouter",
    search: "Rechercher",
    close: "Fermer",
    ok: "OK",
    yes: "Oui",
    no: "Non",
    loading: "Chargement...",
    success: "Succès",
    error: "Erreur",
    warning: "Avertissement",
    info: "Information",
    copy: "Copier",
    paste: "Coller",
    cut: "Couper",
    refresh: "Actualiser",
    reset: "Réinitialiser",
    clear: "Effacer",
    back: "Retour",
    next: "Suivant",
    previous: "Précédent",
    finish: "Terminer",
    submit: "Soumettre",
    import: "Importer",
    export: "Exporter",
    settings: "Paramètres",
    help: "Aide",
    about: "À propos",
    more: "Plus",
  },

  // Page de connexion
  login: {
    title: "Connexion au coffre-fort",
    password: "Mot de passe",
    passwordPlaceholder: "Entrez le mot de passe du coffre-fort",
    login: "Se connecter",
    createVault: "Créer un coffre-fort",
    openVault: "Ouvrir un autre coffre-fort",
    invalidPassword: "Mot de passe incorrect",
    loginSuccess: "Connexion réussie",
    loginFailed: "Échec de la connexion",
    recentVaults: "Récemment utilisé",
    passwordHint1:
      "Le mot de passe de connexion est la clé pour chiffrer et déchiffrer le coffre-fort",
    passwordHint2:
      "Veuillez conserver votre mot de passe en toute sécurité, une fois perdu, il ne peut pas être récupéré",
    passwordHint3:
      "Le mot de passe de connexion est stocké sous forme chiffrée dans la base de données",
    vaultFile: "Fichier du coffre-fort",
    vaultFilePlaceholder: "Veuillez sélectionner le fichier du coffre-fort",
    vaultPathRequired: "Veuillez entrer le chemin du fichier du coffre-fort",
    passwordRequired: "Veuillez entrer le mot de passe de connexion",
    passwordMinLength:
      "La longueur du mot de passe ne peut pas être inférieure à 6 caractères",
    selectFileFailed: "Échec de la sélection du fichier",
    passwordEmpty: "Le mot de passe ne peut pas être vide",
    passwordMinLength8:
      "La longueur du mot de passe doit être d'au moins 8 caractères",
    passwordNeedUppercase:
      "Le mot de passe doit contenir des lettres majuscules",
    passwordNeedLowercase:
      "Le mot de passe doit contenir des lettres minuscules",
    passwordNeedNumber: "Le mot de passe doit contenir des chiffres",
    passwordNeedSpecialChar:
      "Le mot de passe doit contenir des caractères spéciaux (!@#$%^&* etc.)",
    passwordStrengthOk: "La force du mot de passe répond aux exigences",
    vaultFileNotExists: "Le fichier du coffre-fort n'existe pas",
    createVaultPrompt:
      "Veuillez entrer le nom du coffre-fort (sans extension)\n\nLe système créera automatiquement un répertoire de données dans le répertoire du programme actuel et enregistrera le fichier du coffre-fort dans le répertoire de données, créant automatiquement un fichier .db en fonction du nom.",
    vaultNamePlaceholder: "Exemple : mon_coffre_fort",
    setPasswordPrompt:
      "Veuillez définir le mot de passe de connexion\n\nExigences du mot de passe :\n• Au moins 8 caractères\n• Contenir des lettres majuscules\n• Contenir des lettres minuscules\n• Contenir des chiffres\n• Contenir des caractères spéciaux (!@#$%^&* etc.)",
    setPassword: "Définir le mot de passe",
    passwordRequirements: "Exigences du mot de passe",
    passwordReq1: "Au moins 8 caractères",
    passwordReq2: "Contenir des lettres majuscules",
    passwordReq3: "Contenir des lettres minuscules",
    passwordReq4: "Contenir des chiffres",
    passwordReq5: "Contenir des caractères spéciaux (!@#$%^&* etc.)",
    createVaultSuccess: "Coffre-fort créé avec succès",
    createVaultFailed: "Échec de la création du coffre-fort",
    switchToFullMode:
      "Passé en mode complet, veuillez sélectionner le fichier du coffre-fort",
    switchToCreateMode:
      "Passé en mode complet, veuillez créer un nouveau coffre-fort",
  },

  // Interface principale
  main: {
    title: "wepass - Gestionnaire de mots de passe",
    searchPlaceholder: "Rechercher des comptes...",
    noData: "Aucune donnée",
    noSearchResults: "Aucun résultat correspondant trouvé",
    searchResults: "Résultats de la recherche",
    groups: "Groupes",
    group: "Groupe",
    accounts: "Comptes",
    allAccounts: "Tous les comptes",
    noAddress: "Aucune adresse",
    belongsToGroup: "Appartient au groupe",
    inputUsernameAndPassword: "Saisir nom d'utilisateur et mot de passe",
    inputUsername: "Saisir nom d'utilisateur",
    inputPassword: "Saisir mot de passe",
    tryOtherKeywords: "Essayez d'autres mots-clés",
    unknownGroup: "Groupe inconnu",
  },

  // Dialogue des paramètres
  settings: {
    title: "Paramètres",
    general: "Général",
    log: "Journal",
    lock: "Verrouillage",
    theme: "Thème",
    language: "Langue",
    lightTheme: "Clair",
    darkTheme: "Sombre",
    selectTheme: "Sélectionner le thème",
    selectLanguage: "Sélectionner la langue",
    logSettings: "Paramètres du journal",
    enableInfoLog: "Activer le journal d'information",
    enableDebugLog: "Activer le journal de débogage",
    lockSettings: "Paramètres de verrouillage",
    enableAutoLock: "Activer le verrouillage automatique",
    enableTimerLock: "Activer le verrouillage par minuterie",
    enableMinimizeLock: "Activer le verrouillage lors de la minimisation",
    lockTimeMinutes: "Temps de verrouillage (minutes)",
    enableSystemLock: "Activer le verrouillage système",
    systemLockMinutes: "Temps de verrouillage système (minutes)",
    settingsSaved: "Paramètres enregistrés",
    loadSettingsFailed: "Échec du chargement des paramètres",
    saveSettingsFailed: "Échec de l'enregistrement des paramètres",
  },

  // Gestion des comptes
  account: {
    title: "Titre",
    username: "Nom d'utilisateur",
    password: "Mot de passe",
    url: "URL",
    type: "Type",
    notes: "Notes",
    icon: "Icône",
    favorite: "Favori",
    useCount: "Nombre d'utilisations",
    lastUsed: "Dernière utilisation",
    created: "Créé",
    updated: "Mis à jour",
    inputMethod: "Méthode de saisie",
    addAccount: "Ajouter un compte",
    editAccount: "Modifier le compte",
    deleteAccount: "Supprimer le compte",
    copyUsername: "Copier le nom d'utilisateur",
    copyPassword: "Copier le mot de passe",
    openUrl: "Ouvrir l'URL",
    generatePassword: "Générer un mot de passe",
    showPassword: "Afficher le mot de passe",
    hidePassword: "Masquer le mot de passe",
    accountSaved: "Compte enregistré",
    accountDeleted: "Compte supprimé",
    confirmDelete: "Êtes-vous sûr de vouloir supprimer ce compte ?",
    deleteConfirmTitle: "Confirmation de suppression",
    copiedToClipboard: "Copié dans le presse-papiers",
    copyFailed: "Échec de la copie",
    noAccounts: "Aucun compte",
    noAccountsHint: "Cliquez sur le bouton '+' pour créer votre premier compte",
  },

  // Groupes et types
  group: {
    allGroups: "Tous les groupes",
    addGroup: "Ajouter un groupe",
    editGroup: "Modifier le groupe",
    deleteGroup: "Supprimer le groupe",
    groupName: "Nom du groupe",
    groupIcon: "Icône du groupe",
    rename: "Renommer", // TODO: Translate
    createGroup: "Créer un groupe", // TODO: Translate
    moveLeft: "Déplacer à gauche", // TODO: Translate
    moveRight: "Déplacer à droite", // TODO: Translate
    newGroupName: "Nouveau nom de groupe", // TODO: Translate
    renameGroup: "Renommer le groupe", // TODO: Translate
    newGroupPlaceholder: "Nom du groupe", // TODO: Translate
    groupNameCannotBeEmpty: "Le nom du groupe ne peut pas être vide", // TODO: Translate
    groupNameNotChanged: "Le nom du groupe n'a pas changé", // TODO: Translate
    defaultGroupCannotBeDeleted:
      "Le groupe par défaut ne peut pas être supprimé", // TODO: Translate
    confirmDeleteGroup:
      'Êtes-vous sûr de vouloir supprimer le groupe "{groupName}" ?\nToutes les étiquettes et tous les comptes de ce groupe seront également supprimés, cette opération est irréversible.', // TODO: Translate
    deleteGroupTitle: "Supprimer le groupe", // TODO: Translate
    confirmDelete: "Confirmer la suppression", // TODO: Translate
    defaultGroupName: "Défaut", // TODO: Translate
  },

  type: {
    allTypes: "Tous les types",
    addType: "Ajouter un type",
    editType: "Modifier le type",
    deleteType: "Supprimer le type",
    typeName: "Nom du type",
    typeIcon: "Icône du type",
  },

  // Dialogue d'aide
  help: {
    title: "Aide",
    quickStart: "Démarrage rapide",
    welcome: "Bienvenue dans WePassword",
    welcomeTitle: "Bienvenue dans WePassword",
    description:
      "WePassword - est un outil sécurisé de gestion de mots de passe qui vous aide à stocker et gérer vos mots de passe en toute sécurité.",
    welcomeDescription:
      "WePassword - est un outil sécurisé de gestion de mots de passe qui vous aide à stocker et gérer vos mots de passe en toute sécurité.",
    basicOperations: "Opérations de base",
    createVault:
      "Créer un coffre-fort : Lors de la première utilisation, le système vous guidera pour créer un nouveau coffre-fort",
    addPassword:
      'Ajouter un mot de passe : Cliquez sur le bouton "Ajouter" pour saisir les informations du site web, nom d\'utilisateur et mot de passe',
    searchPassword:
      "Rechercher un mot de passe : Utilisez la zone de recherche pour trouver rapidement des entrées de mot de passe spécifiques",
    editPassword:
      "Modifier un mot de passe : Double-cliquez sur une entrée de mot de passe ou cliquez sur le bouton modifier pour apporter des modifications",
    deletePassword:
      "Supprimer un mot de passe : Sélectionnez une entrée de mot de passe et cliquez sur le bouton supprimer",
    features: "Fonctionnalités", // TODO: Translate
    mainFeatures: "Fonctionnalités principales", // TODO: Translate
    passwordGenerator: "Générateur de mots de passe", // TODO: Translate
    passwordGeneratorDesc:
      "Générez des mots de passe forts, prenez en charge les règles personnalisées", // TODO: Translate
    groupManagement: "Gestion de groupe", // TODO: Translate
    groupManagementDesc:
      "Utilisez les onglets pour classer et gérer les mots de passe", // TODO: Translate
    secureEncryption: "Chiffrement sécurisé", // TODO: Translate
    secureEncryptionDesc:
      "Toutes les données de mot de passe sont stockées chiffrées", // TODO: Translate
    importExport: "Importer/Exporter", // TODO: Translate
    importExportDesc:
      "Prend en charge l'importation de données à partir d'autres gestionnaires de mots de passe", // TODO: Translate
    backupRestore: "Sauvegarde/Restauration", // TODO: Translate
    backupRestoreDesc:
      "Sauvegardez régulièrement votre coffre-fort de mots de passe", // TODO: Translate
    passwordGenerationRules: "Règles de génération de mot de passe", // TODO: Translate
    passwordGenerationRulesDesc:
      "Prend en charge plusieurs jeux de caractères et règles personnalisées :", // TODO: Translate
    lowercaseLetters: "Lettres minuscules", // TODO: Translate
    mixedCaseLetters: "Lettres majuscules et minuscules", // TODO: Translate
    uppercaseLetters: "Lettres majuscules", // TODO: Translate
    digits: "Chiffres", // TODO: Translate
    specialCharacters: "Caractères spéciaux", // TODO: Translate
    customCharacterSet: "Jeu de caractères personnalisé", // TODO: Translate
    securityTips: "Conseils de sécurité", // TODO: Translate
    securitySuggestions: "Suggestions de sécurité", // TODO: Translate
    masterPassword: "Mot de passe maître", // TODO: Translate
    masterPasswordDesc:
      "Définissez un mot de passe maître fort et mémorisez-le", // TODO: Translate
    regularBackup: "Sauvegarde régulière", // TODO: Translate
    regularBackupDesc:
      "Sauvegardez régulièrement votre fichier de coffre-fort de mots de passe", // TODO: Translate
    timelyUpdate: "Mise à jour opportune", // TODO: Translate
    timelyUpdateDesc:
      "Mettez à jour régulièrement les mots de passe des comptes importants", // TODO: Translate
    avoidRepetition: "Évitez la répétition", // TODO: Translate
    avoidRepetitionDesc:
      "N'utilisez pas le même mot de passe sur plusieurs sites Web", // TODO: Translate
    safeEnvironment: "Environnement sûr", // TODO: Translate
    safeEnvironmentDesc:
      "Utilisez le gestionnaire de mots de passe dans un environnement sûr", // TODO: Translate
    precautions: "Précautions", // TODO: Translate
    precaution1:
      "Veuillez conserver votre mot de passe maître en lieu sûr. Si vous l'oubliez, vous ne pourrez pas récupérer vos données.", // TODO: Translate
    precaution2:
      "Il est recommandé de changer régulièrement le mot de passe maître pour améliorer la sécurité.", // TODO: Translate
    precaution3:
      "Ne sauvegardez pas le fichier du coffre-fort de mots de passe sur un ordinateur public.", // TODO: Translate
    faq: "FAQ", // TODO: Translate
    faqTitle: "Foire aux questions", // TODO: Translate
    faq1_q: "Que faire si j'oublie mon mot de passe maître ?", // TODO: Translate
    faq1_a:
      "Malheureusement, si vous oubliez votre mot de passe maître, vous ne pourrez pas récupérer les données du coffre-fort de mots de passe. Il est recommandé de sauvegarder régulièrement votre coffre-fort de mots de passe et de noter votre mot de passe maître dans un endroit sûr.", // TODO: Translate
    faq2_q: "Comment sauvegarder le coffre-fort de mots de passe ?", // TODO: Translate
    faq2_a:
      'Vous pouvez copier le fichier du coffre-fort de mots de passe dans un emplacement sûr ou utiliser la fonction d\'exportation dans le menu "Plus".', // TODO: Translate
    faq3_q: "Quels formats de fichiers sont pris en charge ?", // TODO: Translate
    faq3_a:
      "Actuellement, les fichiers de coffre-fort de mots de passe au format .db et .vault sont pris en charge.", // TODO: Translate
    faq4_q: "Comment générer un mot de passe sécurisé ?", // TODO: Translate
    faq4_a:
      'Utilisez la fonction "Générer un mot de passe" dans le menu "Plus". Il est recommandé d\'utiliser une combinaison de lettres majuscules et minuscules, de chiffres et de caractères spéciaux, d\'une longueur d\'au moins 12 caractères.', // TODO: Translate
  },

  // Dialogue À propos
  about: {
    title: "À propos",
    appName: "Gestionnaire de mots de passe",
    description:
      "Outil de gestion de mots de passe multiplateforme basé sur Wails + Go + Vue.js",
    version: "Version",
    buildDate: "Date de compilation",
    author: "Auteur",
    support: "Support",
    license: "Licence",
    github: "GitHub",
    gitee: "Gitee",
  },

  // Messages d'erreur
  error: {
    networkError: "Erreur réseau",
    serverError: "Erreur serveur",
    unknownError: "Erreur inconnue",
    operationFailed: "Opération échouée",
    dataLoadFailed: "Échec du chargement des données",
    saveFailed: "Échec de l'enregistrement",
    deleteFailed: "Échec de la suppression",
    copyFailed: "Échec de la copie",
    invalidInput: "Saisie invalide",
    requiredField: "Ce champ est obligatoire",
    apiServiceUnavailable: "Service API indisponible",
    backendServiceUnavailable:
      "Service backend indisponible, veuillez vérifier si l'application fonctionne correctement",
    accountDataFormatError: "Erreur de format des données de compte",
    accessibilityPermissionRequired:
      "Permission d'accessibilité requise pour la saisie automatique",
    autofillFailed: "Échec du remplissage automatique",
    autofillUsernameFailed:
      "Échec du remplissage automatique du nom d'utilisateur",
    autofillPasswordFailed: "Échec du remplissage automatique du mot de passe",
  },

  // Messages de succès
  success: {
    operationSuccess: "Opération réussie",
    dataSaved: "Données enregistrées",
    dataDeleted: "Données supprimées",
    copied: "Copié",
    imported: "Importation réussie",
    exported: "Exportation réussie",
    dataRefreshed: "Données actualisées",
    tabRenamed: "Onglet renommé avec succès",
    tabDeleted: "Onglet supprimé avec succès",
    tabMoved: "Onglet déplacé avec succès",
    tabCreated: "Onglet créé avec succès",
    groupCreated: "Groupe créé avec succès",
    groupRenamed: "Groupe renommé avec succès",
    groupDeleted: "Groupe supprimé avec succès",
    groupMoved: "Groupe déplacé avec succès",
    accountSaved: "Compte enregistré avec succès",
    accountDeleted: "Compte supprimé avec succès",
    passwordGenerated: "Mot de passe généré avec succès",
    settingsSaved: "Paramètres enregistrés avec succès",
    vaultExported: "Coffre-fort exporté avec succès",
    vaultImported: "Coffre-fort importé avec succès",
    autofillUsernameAndPassword:
      "Remplissage automatique du nom d'utilisateur et du mot de passe réussi",
    autofillUsername: "Remplissage automatique du nom d'utilisateur réussi",
    autofillPassword: "Remplissage automatique du mot de passe réussi",
    usernameCopied:
      "Nom d'utilisateur copié dans le presse-papiers (effacement automatique après 10 secondes)", // TODO: Translate
    passwordCopied: "Mot de passe copié dans le presse-papiers", // TODO: Translate
    usernameAndPasswordCopied:
      "Nom d'utilisateur et mot de passe copiés dans le presse-papiers (effacement automatique après 10 secondes)", // TODO: Translate
    urlCopied: "URL copiée dans le presse-papiers", // TODO: Translate
    titleCopied: "Titre copié dans le presse-papiers", // TODO: Translate
    notesCopied:
      "Notes copiées dans le presse-papiers (effacement automatique après 10 secondes)", // TODO: Translate
  },

  // Messages d'avertissement
  warning: {
    noGroupData: "Aucune donnée de groupe, veuillez d'abord créer un groupe",
    selectGroupFirst: "Veuillez d'abord sélectionner un groupe",
    selectTabFirst: "Veuillez d'abord sélectionner un onglet",
    noAccountSelected: "Veuillez d'abord sélectionner un compte",
    confirmOperation: "Veuillez confirmer cette opération",
  },

  // Informations de statut
  status: {
    loading: "Chargement...",
    saving: "Enregistrement...",
    deleting: "Suppression...",
    processing: "Traitement...",
    connecting: "Connexion...",
    searching: "Recherche...",
    exporting: "Exportation...",
    importing: "Importation...",
    generating: "Génération...",
  },

  // Menu Plus
  moreMenu: {
    more: "Plus",
    selectNewVault: "Sélectionner un nouveau coffre-fort",
    openVaultDirectory: "Ouvrir le répertoire du coffre-fort",
    generatePassword: "Générer un mot de passe",
    setPasswordRules: "Définir les règles de génération de mot de passe",
    changeLoginPassword: "Modifier le mot de passe de connexion",
    exportVault: "Exporter le coffre-fort",
    importVault: "Importer le coffre-fort",
    changeLog: "Journal des modifications",
    settings: "Paramètres",
    help: "Aide",
    about: "À propos",
    lockVault: "Verrouiller le coffre-fort",
    logout: "Se déconnecter",
    oldLoginPasswordLabel: "Ancien mot de passe",
    oldLoginPasswordPlaceholder: "Entrez le mot de passe de connexion actuel",
    newLoginPasswordLabel: "Nouveau mot de passe",
    newLoginPasswordPlaceholder: "Entrez le nouveau mot de passe de connexion",
    confirmNewPasswordLabel: "Confirmer le nouveau mot de passe",
    confirmNewPasswordPlaceholder:
      "Entrez à nouveau le nouveau mot de passe de connexion",
    selectVaultFile: "Sélectionner le fichier du coffre-fort",
    selectVaultFilePrompt:
      "Veuillez sélectionner le fichier du coffre-fort à ouvrir :",
    selectVaultFilePlaceholder:
      "Veuillez sélectionner le chemin du fichier du coffre-fort",
    browse: "Parcourir",
    supportedFileFormats: "Formats de fichier pris en charge : .db, .vault",
    openFileFailed: "Échec de l'ouverture du fichier",
    passwordCopied: "Mot de passe copié dans le presse-papiers",
    copyFailedManual:
      "Échec de la copie, veuillez copier le mot de passe manuellement",
    passwordChangeInProgress:
      "Le changement de mot de passe est en cours, veuillez patienter...",
    formRefNotFound: "Référence du formulaire introuvable",
    oldPasswordIncorrect:
      "Ancien mot de passe incorrect, veuillez le saisir à nouveau",
    changePasswordConfirm:
      "Changer le mot de passe de connexion rechiffrera toutes les données du compte. Cette opération est irréversible, assurez-vous d'avoir sauvegardé votre coffre-fort actuel. Voulez-vous continuer ?",
    passwordChangeSuccess: "Mot de passe de connexion changé avec succès !",
    passwordChangeFailed: "Échec du changement de mot de passe de connexion",
    importSuccessDataUpdated: "Importation réussie, données mises à jour",
    importSuccessDataRefreshFailed:
      "Importation réussie, mais l'actualisation des données a échoué, veuillez actualiser la page manuellement",
    vaultDirectoryOpened: "Répertoire du coffre-fort ouvert",
    openDirectoryFailed: "Échec de l'ouverture du répertoire",
    lockVaultConfirm:
      "Êtes-vous sûr de vouloir verrouiller le coffre-fort ? Après le verrouillage, vous devrez entrer à nouveau le mot de passe pour y accéder.",
    vaultLocked: "Coffre-fort verrouillé",
    lockVaultFailed: "Échec du verrouillage du coffre-fort",
    logoutConfirm:
      "Êtes-vous sûr de vouloir vous déconnecter ? Après la déconnexion, vous devrez entrer à nouveau le mot de passe pour accéder au coffre-fort.",
    logoutSuccess: "Déconnexion réussie",
    oldLoginPasswordRequired: "Entrez l'ancien mot de passe",
    newLoginPasswordRequired: "Entrez le nouveau mot de passe",
    newLoginPasswordMinLength:
      "Le mot de passe doit contenir au moins 8 caractères",
    newLoginPasswordStrength:
      "Le mot de passe doit contenir des majuscules, des minuscules, des chiffres et des caractères spéciaux",
    confirmNewLoginPasswordRequired: "Confirmez le nouveau mot de passe",
    passwordsDoNotMatch: "Les mots de passe ne correspondent pas",
    selectNewVaultConfirm:
      "Sélectionner un nouveau coffre-fort vous ramènera à l'écran de connexion, où vous devrez sélectionner à nouveau le fichier du coffre-fort et entrer le mot de passe. Voulez-vous continuer ?",
    selectNewVaultSuccess:
      "Vous êtes revenu à l'écran de connexion, sélectionnez un nouveau coffre-fort",
    pleaseSelectFile: "Veuillez d'abord sélectionner un fichier",
    openingFile: "Ouverture du fichier : ",
    open: "Ouvrir",
    continue: "Continuer",
    verifying: "Vérification...",
    changingPassword: "Changement...",
  },

  // Menu contextuel
  contextMenu: {
    inputUsernameAndPassword: "Saisir nom d'utilisateur et mot de passe", // TODO: Translate
    openUrl: "Ouvrir l'adresse",
    duplicate: "Créer une copie",
    view: "Voir",
    edit: "Modifier",
    changeGroup: "Changer de groupe",
    copyUsername: "Copier le nom d'utilisateur",
    copyPassword: "Copier le mot de passe",
    copyUsernameAndPassword: "Copier nom d'utilisateur et mot de passe",
    copyUrl: "Copier l'adresse",
    copyTitle: "Copier le titre",
    copyNotes: "Copier les notes",
    delete: "Supprimer",
  },

  // Fonction d'exportation
  export: {
    title: "Exporter le coffre-fort",
    steps: {
      verifyPassword: "Vérifier le mot de passe",
      selectAccounts: "Sélectionner les comptes",
      setBackup: "Configurer la sauvegarde",
      exportComplete: "Exportation terminée",
    },
    verifyPasswordTitle: "Vérifier le mot de passe de connexion",
    verifyPasswordDesc:
      "Veuillez entrer le mot de passe de connexion actuel du coffre-fort pour continuer l'opération d'exportation.",
    loginPassword: "Mot de passe de connexion",
    loginPasswordPlaceholder: "Veuillez entrer le mot de passe de connexion",
    selectAccountsTitle: "Sélectionner les comptes à exporter",
    exportAll: "Tout exporter",
    exportByGroup: "Exporter par groupe",
    exportByType: "Exporter par type", // TODO: Translate
    exportSelected: "Exporter la sélection",
    selectGroups: "Sélectionner les groupes", // TODO: Translate
    selectAll: "Tout sélectionner", // TODO: Translate
    clearAll: "Tout désélectionner", // TODO: Translate
    groupSelectionSummary:
      "{count} groupes sélectionnés, {accountCount} comptes attendus", // TODO: Translate
    selectTypes: "Sélectionner les types", // TODO: Translate
    typeSelectionSummary:
      "{count} types sélectionnés, {accountCount} comptes attendus", // TODO: Translate
    selectAccounts: "Sélectionner les comptes", // TODO: Translate
    loadingAccounts: "Chargement des comptes...", // TODO: Translate
    noAccounts: "Aucune donnée de compte", // TODO: Translate
    setBackupTitle:
      "Définir le mot de passe de sauvegarde et le chemin d'exportation", // TODO: Translate
    backupPassword: "Mot de passe de sauvegarde", // TODO: Translate
    backupPasswordPlaceholder: "Veuillez entrer le mot de passe de sauvegarde", // TODO: Translate
    generate: "Générer", // TODO: Translate
    backupPasswordTip:
      "Le mot de passe de sauvegarde est utilisé pour chiffrer les données exportées, veuillez le conserver en lieu sûr", // TODO: Translate
    exportPath: "Chemin d'exportation", // TODO: Translate
    exportPathPlaceholder: "Veuillez sélectionner le chemin d'exportation", // TODO: Translate
    browse: "Parcourir", // TODO: Translate
    exporting: "Exportation du coffre-fort...", // TODO: Translate
    exportSuccessTitle: "Exportation réussie", // TODO: Translate
    exportSuccessSubTitle:
      "Le coffre-fort a été exporté avec succès vers : {path}", // TODO: Translate
    exportFailedTitle: "Échec de l'exportation", // TODO: Translate
    startExport: "Démarrer l'exportation", // TODO: Translate
    openFolder: "Ouvrir le dossier", // TODO: Translate
    loginPasswordRequired: "Veuillez entrer le mot de passe de connexion", // TODO: Translate
    backupPasswordRequired: "Veuillez entrer le mot de passe de sauvegarde", // TODO: Translate
    backupPasswordMinLength:
      "Le mot de passe de sauvegarde doit comporter au moins 6 caractères", // TODO: Translate
    exportPathRequired: "Veuillez sélectionner le chemin d'exportation", // TODO: Translate
  },

  // Importar
  import: {
    title: "Importer le coffre-fort",
    selectFile: "Sélectionner un fichier",
    selectFileDesc:
      "Veuillez sélectionner le fichier du coffre-fort à importer",
    fileFormat: "Format de fichier",
    importProgress: "Progression de l'importation",
    importComplete: "Importation terminée",
  },

  // 新增的导入功能翻译
  importVault: {
    title: "Importer le coffre-fort de mots de passe",
    step1: "Sélectionner un fichier",
    step2: "Vérifier le mot de passe",
    step3: "Importation terminée",
    step1Title: "Sélectionner le fichier d'importation",
    step1Description:
      "Veuillez sélectionner le fichier de sauvegarde du coffre-fort de mots de passe (format ZIP) à importer.",
    importFile: "Fichier d'importation",
    selectImportFilePlaceholder:
      "Veuillez sélectionner le fichier d'importation",
    browse: "Parcourir",
    step2Title: "Entrer le mot de passe de décompression",
    step2Description:
      "Veuillez entrer le mot de passe de décompression (mot de passe de sauvegarde) du fichier.",
    backupPassword: "Mot de passe de décompression",
    enterBackupPasswordPlaceholder:
      "Veuillez entrer le mot de passe de décompression",
    importingVault: "Importation du coffre-fort de mots de passe...",
    importComplete: "Importation terminée",
    importSuccess: "Importation réussie",
    vaultImportedSuccessfully:
      "Le coffre-fort de mots de passe a été importé avec succès",
    importReport: "Rapport d'importation",
    totalAccounts: "Nombre total de comptes",
    successfullyImported: "Importés avec succès",
    skippedAccounts: "Comptes ignorés",
    errorAccounts: "Comptes en erreur",
    totalGroups: "Nombre total de groupes",
    importedGroups: "Groupes importés",
    totalTypes: "Nombre total de types",
    importedTypes: "Types importés",
    skippedAccountDetails: "Détails des comptes ignorés",
    accountTitle: "Titre du compte",
    accountName: "Nom du compte",
    accountId: "ID du compte",
    importFailed: "Échec de l'importation",
    importFailedStatus: "Échec de l'importation",
    unknownError: "Une erreur inconnue s'est produite pendant l'importation",
    close: "Fermer",
    cancel: "Annuler",
    previousStep: "Étape précédente",
    startImport: "Démarrer l'importation",
    nextStep: "Étape suivante",
    refreshData: "Actualiser les données",
    selectImportFileMessage: "Veuillez sélectionner le fichier d'importation",
    enterBackupPasswordMessage:
      "Veuillez entrer le mot de passe de décompression",
    selectFileSuccess: "Fichier d'importation sélectionné avec succès",
    selectFileFailed: "Échec de la sélection du fichier d'importation",
    preparingToImport: "Préparation de l'importation...",
    validatingFileAndPassword: "Validation du fichier et du mot de passe...",
    vaultImportSuccess: "Coffre-fort de mots de passe importé avec succès",
    dataRefreshed: "Données actualisées",
    importInProgressWarning:
      "L'importation est en cours, veuillez patienter...",
  },

  // Générateur de mot de passe
  passwordGenerator: {
    title: "Générateur de mots de passe",
    selectRule: "Sélectionner une règle de mot de passe",
    selectRulePlaceholder: "Veuillez sélectionner une règle de mot de passe",
    generalRule: "Règle de mot de passe générale",
    customRule: "Règle de mot de passe personnalisée",
    includeUppercase: "Inclure les lettres majuscules",
    includeLowercase: "Inclure les lettres minuscules",
    includeNumbers: "Inclure les chiffres",
    includeSpecialChars: "Inclure les caractères spéciaux",
    passwordLength: "Longueur du mot de passe",
    customSpecialChars: "Caractères spéciaux personnalisés",
    defaultSpecialChars:
      "Laisser vide pour utiliser les caractères spéciaux par défaut",
    generatePassword: "Générer un mot de passe",
    usePassword: "Utiliser ce mot de passe",
    generatedPassword: "Mot de passe généré",
    clickToGenerate: "Cliquez sur le bouton générer un mot de passe",
    selectCharType: "Veuillez sélectionner au moins un type de caractère",
    enterPattern: "Veuillez entrer un modèle de mot de passe",
    generateFirst: "Veuillez d'abord générer un mot de passe",
  },

  // Nouvelles traductions pour les paramètres des règles de mot de passe
  passwordRuleSettings: {
    title: "Paramètres des règles de génération de mot de passe",
    savedRules: "Règles de mot de passe enregistrées",
    newRule: "Nouvelle règle",
    ruleName: "Nom de la règle",
    description: "Description",
    operation: "Opération",
    edit: "Modifier",
    delete: "Supprimer",
    generalRule: "Règle générale de mot de passe",
    includeUppercase: "Inclure les majuscules",
    includeLowercase: "Inclure les minuscules",
    includeNumbers: "Inclure les chiffres",
    includeSpecialChars: "Inclure les caractères spéciaux",
    passwordLength: "Longueur du mot de passe",
    customSpecialChars: "Caractères spéciaux personnalisés",
    customSpecialCharsPlaceholder:
      "Laisser vide pour utiliser les caractères spéciaux par défaut : !@#$%^&*()_+-=[]{}|;:,.<>?",
    customRulesDescription:
      "Description des règles de mot de passe personnalisées",
    close: "Fermer",
    saveSettings: "Enregistrer les paramètres",
    editRule: "Modifier la règle",
    newRuleDialog: "Nouvelle règle",
    ruleNameRequired: "Veuillez entrer le nom de la règle",
    ruleDescriptionPlaceholder: "Veuillez entrer la description de la règle",
    ruleType: "Type de règle",
    general: "Règle générale",
    custom: "Règle personnalisée",
    passwordPattern: "Modèle de mot de passe",
    passwordPatternPlaceholder: "Exemple : Aaa111",
    save: "Enregistrer",
    cancel: "Annuler",
    confirmDeleteRule:
      'Êtes-vous sûr de vouloir supprimer la règle "{ruleName}"?',
    deleteConfirmation: "Confirmation de suppression",
    confirmDelete: "Confirmer",
    ruleDeletedSuccess: "Règle supprimée avec succès",
    ruleUpdateSuccess: "Règle mise à jour avec succès",
    ruleCreateSuccess: "Règle créée avec succès",
    enterRuleName: "Veuillez entrer le nom de la règle",
    settingsSaved: "Paramètres enregistrés avec succès",
  },

  // Nouvelles traductions pour les résultats de recherche
  searchResults: {
    groups: "Groupes",
    accounts: "Comptes",
    noUrl: "Aucune URL",
    belongsToGroup: "Appartient au groupe",
    unknownGroup: "Groupe inconnu",
    noSearchResults: "Aucun résultat de recherche trouvé",
    tryOtherKeywords: "Essayez d'autres mots-clés",
  },

  // Nouvelles traductions pour la barre d'état
  statusBar: {
    authInfo: "Informations d'autorisation : Version Pro activée",
    usageCount: "Utilisations : {count} fois",
    usageDays: "Jours d'utilisation : {days} jours",
  },

  // Nouvelles traductions pour le menu contextuel des onglets
  tabContextMenu: {
    rename: "Renommer",
    deleteTab: "Supprimer l'onglet",
    newTab: "Nouvel onglet",
    moveUp: "Déplacer vers le haut",
    moveDown: "Déplacer vers le bas",
    promptNewName: "Veuillez entrer le nouveau nom de l'onglet",
    renameTab: "Renommer l'onglet",
    confirm: "Confirmer",
    cancel: "Annuler",
    tabName: "Nom de l'onglet",
    tabNameCannotBeEmpty: "Le nom de l'onglet ne peut pas être vide",
    tabNameNotChanged: "Le nom de l'onglet n'a pas changé",
    userCanceledRename: "L'utilisateur a annulé l'opération de renommage",
    confirmDeleteTab:
      'Êtes-vous sûr de vouloir supprimer l\'onglet "{tabName}"?\nAprès la suppression, tous les comptes sous cet onglet seront également supprimés, cette opération est irréversible.',
    deleteConfirmation: "Supprimer l'onglet",
    confirmDelete: "Confirmer la suppression",
    userCanceledDelete: "L'utilisateur a annulé l'opération de suppression",
    promptNewTabName: "Veuillez entrer le nouveau nom de l'onglet",
    userCanceledNewTab: "L'utilisateur a annulé l'opération de nouvel onglet",
  },

  // Nouvelles traductions pour la barre latérale des onglets
  tabsSidebar: {
    emptyTabs: "Aucun onglet",
    createTabHint: "Cliquez sur le bouton ci-dessous pour créer un onglet",
    newTab: "Nouvel onglet",
  },

  // Nouvelles traductions pour la boîte de dialogue de test
  testDialog: {
    title: "Boîte de dialogue de test",
    description1:
      "Ceci est une boîte de dialogue de test pour vérifier si la fonction de fenêtre contextuelle fonctionne correctement.",
    description2:
      "Si vous pouvez voir cette boîte de dialogue, cela signifie que la fonction de fenêtre contextuelle fonctionne correctement.",
    close: "Fermer",
  },

  // Nouvelles traductions pour la barre de titre
  titleBar: {
    defaultAppTitle: "Gestionnaire de mots de passe",
  },

  // 新增的锁定事件服务翻译
  lockEventService: {
    logPrefix: "Service d'événements de verrouillage",
    frontendStateUpdated: "État du frontend mis à jour",
    redirectToLogin: "Redirection vers la page de connexion",
    vaultAutoLocked:
      "Le coffre-fort de mots de passe a été automatiquement verrouillé, veuillez vous reconnecter",
    sensitiveDataCleared: "Données sensibles effacées",
    clearSensitiveDataFailed: "Échec de l'effacement des données sensibles",
    windowEventListenersSet:
      "Les écouteurs d'événements de fenêtre ont été définis",
    windowLostFocus: "La fenêtre a perdu le focus",
    windowGainedFocus: "La fenêtre a gagné le focus",
    backendNotifiedMinimize:
      "Le backend a été informé de la minimisation de la fenêtre",
    notifyMinimizeFailed:
      "Échec de la notification de minimisation de la fenêtre",
    backendNotifiedFocus:
      "Le backend a été informé que la fenêtre a gagné le focus",
    notifyFocusFailed:
      "Échec de la notification que la fenêtre a gagné le focus",
    manualLockCheckTriggered:
      "Vérification manuelle du verrouillage déclenchée",
  },

  // Nouvelles traductions pour les fonctions utilitaires de compte
  accountUtils: {
    defaultAccountTitle: "Compte",
    startCopyingPassword:
      "Début de la copie du mot de passe, ID du compte : {accountId}",
    passwordCopiedSuccess:
      "Mot de passe copié dans le presse-papiers (nettoyage automatique après 10 secondes)",
    passwordCopySuccess:
      "Copie du mot de passe réussie, ID du compte : {accountId}",
    passwordCopyFailed: "Échec de la copie du mot de passe",
    startGettingPassword:
      "Début de la récupération du mot de passe, ID du compte : {accountId}",
    passwordGetSuccess:
      "Récupération du mot de passe réussie, ID du compte : {accountId}",
    passwordGetFailed: "Échec de la récupération du mot de passe",
    startGettingAccountDetail:
      "Début de la récupération des détails du compte, ID du compte : {accountId}",
    accountDetailGetSuccess:
      "Récupération des détails du compte réussie, ID du compte : {accountId}",
    accountDetailGetFailed: "Échec de la récupération des détails du compte",
  },
};
