/**
 * Paquete de idioma español
 * @author Chen Fengqing
 * @date 2025-10-05
 */

export default {
  // Común
  common: {
    confirm: "Confirmar",
    cancel: "Cancelar",
    save: "Guardar",
    delete: "Eliminar",
    edit: "Editar",
    add: "Añadir",
    search: "Buscar",
    close: "Cerrar",
    ok: "Aceptar",
    yes: "Sí",
    no: "No",
    loading: "Cargando...",
    success: "Éxito",
    error: "Error",
    warning: "Advertencia",
    info: "Información",
    copy: "Copiar",
    paste: "Pegar",
    cut: "Cortar",
    refresh: "Actualizar",
    reset: "Restablecer",
    clear: "Limpiar",
    back: "Atrás",
    next: "Siguiente",
    previous: "Anterior",
    finish: "Finalizar",
    submit: "Enviar",
    import: "Importar",
    export: "Exportar",
    settings: "Configuración",
    help: "Ayuda",
    about: "Acerca de",
    more: "Más",
  },

  // Página de inicio de sesión
  login: {
    title: "Inicio de sesión en la bóveda",
    password: "Contraseña",
    passwordPlaceholder: "Ingrese la contraseña de la bóveda",
    login: "Iniciar sesión",
    createVault: "Crear bóveda",
    openVault: "Abrir otra bóveda",
    invalidPassword: "Contraseña incorrecta",
    loginSuccess: "Inicio de sesión exitoso",
    loginFailed: "Error al iniciar sesión",
    recentVaults: "Usado recientemente",
    passwordHint1:
      "La contraseña de inicio de sesión es la clave para cifrar y descifrar la bóveda",
    passwordHint2:
      "Por favor guarde bien su contraseña, una vez perdida no se puede recuperar",
    passwordHint3:
      "La contraseña de inicio de sesión se almacena cifrada en la base de datos",
    vaultFile: "Archivo de bóveda",
    vaultFilePlaceholder: "Por favor seleccione el archivo de bóveda",
    vaultPathRequired: "Por favor ingrese la ruta del archivo de bóveda",
    passwordRequired: "Por favor ingrese la contraseña de inicio de sesión",
    passwordMinLength:
      "La longitud de la contraseña no puede ser menor a 6 caracteres",
    selectFileFailed: "Error al seleccionar archivo",
    passwordEmpty: "La contraseña no puede estar vacía",
    passwordMinLength8:
      "La longitud de la contraseña debe ser al menos 8 caracteres",
    passwordNeedUppercase: "La contraseña debe contener letras mayúsculas",
    passwordNeedLowercase: "La contraseña debe contener letras minúsculas",
    passwordNeedNumber: "La contraseña debe contener números",
    passwordNeedSpecialChar:
      "La contraseña debe contener caracteres especiales (!@#$%^&* etc.)",
    passwordStrengthOk: "La fortaleza de la contraseña cumple los requisitos",
    vaultFileNotExists: "El archivo de bóveda no existe",
    createVaultPrompt:
      "Por favor ingrese el nombre de la bóveda (sin extensión)\n\nEl sistema creará automáticamente un directorio data en el directorio actual del programa y guardará el archivo de bóveda en el directorio data, creando automáticamente un archivo .db según el nombre.",
    vaultNamePlaceholder: "Ejemplo: mi_boveda",
    setPasswordPrompt:
      "Por favor configure la contraseña de inicio de sesión\n\nRequisitos de contraseña:\n• Al menos 8 caracteres\n• Contener letras mayúsculas\n• Contener letras minúsculas\n• Contener números\n• Contener caracteres especiales (!@#$%^&* etc.)",
    setPassword: "Configurar contraseña",
    passwordRequirements: "Requisitos de contraseña",
    passwordReq1: "Al menos 8 caracteres",
    passwordReq2: "Contener letras mayúsculas",
    passwordReq3: "Contener letras minúsculas",
    passwordReq4: "Contener números",
    passwordReq5: "Contener caracteres especiales (!@#$%^&* etc.)",
    createVaultSuccess: "Bóveda creada exitosamente",
    createVaultFailed: "Error al crear la bóveda",
    switchToFullMode:
      "Cambiado a modo completo, por favor seleccione el archivo de bóveda",
    switchToCreateMode:
      "Cambiado a modo completo, por favor cree una nueva bóveda",
  },

  // Interfaz principal
  main: {
    title: "wepass",
    searchPlaceholder: "Buscar cuentas...",
    noData: "Sin datos",
    noSearchResults: "No se encontraron resultados coincidentes",
    searchResults: "Resultados de búsqueda",
    groups: "Grupos",
    group: "Grupo",
    accounts: "Cuentas",
    allAccounts: "Todas las cuentas",
    noAddress: "Sin dirección",
    belongsToGroup: "Pertenece al grupo",
    inputUsernameAndPassword: "Introducir nombre de usuario y contraseña",
    inputUsername: "Introducir nombre de usuario",
    inputPassword: "Introducir contraseña",
    tryOtherKeywords: "Intente con otras palabras clave",
    unknownGroup: "Grupo desconocido",
  },

  // Diálogo de configuración
  settings: {
    title: "Configuración",
    general: "General",
    log: "Registro",
    lock: "Bloqueo",
    theme: "Tema",
    language: "Idioma",
    lightTheme: "Claro",
    darkTheme: "Oscuro",
    selectTheme: "Seleccionar tema",
    selectLanguage: "Seleccionar idioma",
    logSettings: "Configuración de registro",
    enableInfoLog: "Habilitar registro de información",
    enableDebugLog: "Habilitar registro de depuración",
    lockSettings: "Configuración de bloqueo",
    enableAutoLock: "Habilitar bloqueo automático",
    enableTimerLock: "Habilitar bloqueo por temporizador",
    enableMinimizeLock: "Habilitar bloqueo al minimizar",
    lockTimeMinutes: "Tiempo de bloqueo (minutos)",
    enableSystemLock: "Habilitar bloqueo del sistema",
    systemLockMinutes: "Tiempo de bloqueo del sistema (minutos)",
    settingsSaved: "Configuración guardada",
    loadSettingsFailed: "Error al cargar la configuración",
    saveSettingsFailed: "Error al guardar la configuración",
  },

  // Gestión de cuentas
  account: {
    title: "Título",
    username: "Nombre de usuario",
    password: "Contraseña",
    url: "URL",
    type: "Tipo",
    notes: "Notas",
    icon: "Icono",
    favorite: "Favorito",
    useCount: "Número de usos",
    lastUsed: "Último uso",
    created: "Creado",
    updated: "Actualizado",
    inputMethod: "Método de entrada",
    addAccount: "Añadir cuenta",
    editAccount: "Editar cuenta",
    deleteAccount: "Eliminar cuenta",
    copyUsername: "Copiar nombre de usuario",
    copyPassword: "Copiar contraseña",
    openUrl: "Abrir URL",
    generatePassword: "Generar contraseña",
    showPassword: "Mostrar contraseña",
    hidePassword: "Ocultar contraseña",
    accountSaved: "Cuenta guardada",
    accountDeleted: "Cuenta eliminada",
    confirmDelete: "¿Está seguro de que desea eliminar esta cuenta?",
    deleteConfirmTitle: "Confirmación de eliminación",
    copiedToClipboard: "Copiado al portapapeles",
    copyFailed: "Error al copiar",
    noAccounts: "Sin cuentas",
    noAccountsHint: "Haga clic en el botón '+' para crear su primera cuenta",
  },

  // Grupos y tipos
  group: {
    allGroups: "Todos los grupos",
    addGroup: "Añadir grupo",
    editGroup: "Editar grupo",
    deleteGroup: "Eliminar grupo",
    groupName: "Nombre del grupo",
    groupIcon: "Icono del grupo",
    rename: "Renombrar", // TODO: Translate
    createGroup: "Crear grupo", // TODO: Translate
    moveLeft: "Mover a la izquierda", // TODO: Translate
    moveRight: "Mover a la derecha", // TODO: Translate
    newGroupName: "Nuevo nombre de grupo", // TODO: Translate
    renameGroup: "Renombrar grupo", // TODO: Translate
    newGroupPlaceholder: "Nombre del grupo", // TODO: Translate
    groupNameCannotBeEmpty: "El nombre del grupo no puede estar vacío", // TODO: Translate
    groupNameNotChanged: "El nombre del grupo no ha cambiado", // TODO: Translate
    defaultGroupCannotBeDeleted: "El grupo predeterminado no se puede eliminar", // TODO: Translate
    confirmDeleteGroup:
      '¿Está seguro de que desea eliminar el grupo "{groupName}"?\nTodas las etiquetas y cuentas de este grupo también se eliminarán, esta operación es irreversible.', // TODO: Translate
    deleteGroupTitle: "Eliminar grupo", // TODO: Translate
    confirmDelete: "Confirmar eliminación", // TODO: Translate
    defaultGroupName: "Predeterminado", // TODO: Translate
  },

  type: {
    allTypes: "Todos los tipos",
    addType: "Añadir tipo",
    editType: "Editar tipo",
    deleteType: "Eliminar tipo",
    typeName: "Nombre del tipo",
    typeIcon: "Icono del tipo",
  },

  // Diálogo de ayuda
  help: {
    title: "Ayuda",
    quickStart: "Inicio rápido",
    welcome: "Bienvenido a WePassword",
    welcomeTitle: "Bienvenido a WePassword",
    description:
      "WePassword es una herramienta segura de gestión de contraseñas que le ayuda a almacenar y gestionar contraseñas de forma segura.",
    welcomeDescription:
      "WePassword es una herramienta segura de gestión de contraseñas que le ayuda a almacenar y gestionar contraseñas de forma segura.",
    basicOperations: "Operaciones básicas",
    createVault:
      "Crear bóveda: Al usar por primera vez, el sistema le guiará para crear una nueva bóveda",
    addPassword:
      'Añadir contraseña: Haga clic en el botón "Añadir" para introducir información del sitio web, nombre de usuario y contraseña',
    searchPassword:
      "Buscar contraseña: Use el cuadro de búsqueda para encontrar rápidamente entradas de contraseña específicas",
    editPassword:
      "Editar contraseña: Haga doble clic en una entrada de contraseña o haga clic en el botón editar para modificar",
    deletePassword:
      "Eliminar contraseña: Seleccione una entrada de contraseña y haga clic en el botón eliminar",
    features: "Características", // TODO: Translate
    mainFeatures: "Características principales", // TODO: Translate
    passwordGenerator: "Generador de contraseñas", // TODO: Translate
    passwordGeneratorDesc:
      "Genere contraseñas seguras, admita reglas personalizadas", // TODO: Translate
    groupManagement: "Gestión de grupos", // TODO: Translate
    groupManagementDesc:
      "Use pestañas para clasificar y administrar contraseñas", // TODO: Translate
    secureEncryption: "Cifrado seguro", // TODO: Translate
    secureEncryptionDesc:
      "Todos los datos de la contraseña se almacenan cifrados", // TODO: Translate
    importExport: "Importar/Exportar", // TODO: Translate
    importExportDesc:
      "Admite la importación de datos de otros administradores de contraseñas", // TODO: Translate
    backupRestore: "Copia de seguridad/Restauración", // TODO: Translate
    backupRestoreDesc:
      "Haga una copia de seguridad de su bóveda de contraseñas con regularidad", // TODO: Translate
    passwordGenerationRules: "Reglas de generación de contraseñas", // TODO: Translate
    passwordGenerationRulesDesc:
      "Admite múltiples juegos de caracteres y reglas personalizadas:", // TODO: Translate
    lowercaseLetters: "Letras minúsculas", // TODO: Translate
    mixedCaseLetters: "Letras mayúsculas y minúsculas", // TODO: Translate
    uppercaseLetters: "Letras mayúsculas", // TODO: Translate
    digits: "Dígitos", // TODO: Translate
    specialCharacters: "Caracteres especiales", // TODO: Translate
    customCharacterSet: "Juego de caracteres personalizado", // TODO: Translate
    securityTips: "Consejos de seguridad", // TODO: Translate
    securitySuggestions: "Sugerencias de seguridad", // TODO: Translate
    masterPassword: "Contraseña maestra", // TODO: Translate
    masterPasswordDesc: "Establezca una contraseña maestra segura y recuérdela", // TODO: Translate
    regularBackup: "Copia de seguridad periódica", // TODO: Translate
    regularBackupDesc:
      "Haga una copia de seguridad de su archivo de bóveda de contraseñas con regularidad", // TODO: Translate
    timelyUpdate: "Actualización oportuna", // TODO: Translate
    timelyUpdateDesc:
      "Actualice periódicamente las contraseñas de las cuentas importantes", // TODO: Translate
    avoidRepetition: "Evite la repetición", // TODO: Translate
    avoidRepetitionDesc: "No use la misma contraseña en varios sitios web", // TODO: Translate
    safeEnvironment: "Entorno seguro", // TODO: Translate
    safeEnvironmentDesc:
      "Use el administrador de contraseñas en un entorno seguro", // TODO: Translate
    precautions: "Precauciones", // TODO: Translate
    precaution1:
      "Guarde su contraseña maestra de forma segura. Si la olvida, no podrá recuperar sus datos.", // TODO: Translate
    precaution2:
      "Se recomienda cambiar la contraseña maestra periódicamente para mejorar la seguridad.", // TODO: Translate
    precaution3:
      "No guarde el archivo de la bóveda de contraseñas en una computadora pública.", // TODO: Translate
    faq: "Preguntas frecuentes", // TODO: Translate
    faqTitle: "Preguntas frecuentes", // TODO: Translate
    faq1_q: "¿Qué hago si olvido mi contraseña maestra?", // TODO: Translate
    faq1_a:
      "Lamentablemente, si olvida su contraseña maestra, no podrá recuperar los datos de la bóveda de contraseñas. Se recomienda que haga una copia de seguridad de su bóveda de contraseñas con regularidad y anote su contraseña maestra en un lugar seguro.", // TODO: Translate
    faq2_q: "¿Cómo hago una copia de seguridad de la bóveda de contraseñas?", // TODO: Translate
    faq2_a:
      'Puede copiar el archivo de la bóveda de contraseñas a una ubicación segura o usar la función de exportación en el menú "Más".', // TODO: Translate
    faq3_q: "¿Qué formatos de archivo son compatibles?", // TODO: Translate
    faq3_a:
      "Actualmente, se admiten archivos de bóveda de contraseñas en formato .db y .vault.", // TODO: Translate
    faq4_q: "¿Cómo genero una contraseña segura?", // TODO: Translate
    faq4_a:
      'Use la función "Generar contraseña" en el menú "Más". Se recomienda usar una combinación de letras mayúsculas y minúsculas, números y caracteres especiales, con una longitud de al menos 12 caracteres.', // TODO: Translate
  },

  // Diálogo Acerca de
  about: {
    title: "Acerca de",
    appName: "Gestor de contraseñas",
    description:
      "Herramienta de gestión de contraseñas multiplataforma basada en Wails + Go + Vue.js",
    version: "Versión",
    buildDate: "Fecha de compilación",
    author: "Autor",
    support: "Soporte",
    license: "Licencia",
    github: "GitHub",
    gitee: "Gitee",
  },

  // 对话框
  dialog: {
    changeGroupTitle: "Cambiar grupo",
    accountTitleLabel: "Título de la cuenta",
    selectGroupLabel: "Seleccionar grupo",
    selectGroupPlaceholder: "Por favor seleccione un grupo",
    selectTypeLabel: "Seleccionar tipo",
    selectTypePlaceholder: "Por favor seleccione un tipo",
    changeLogTitle: "Registro de cambios",
  },

  // Mensajes de error
  error: {
    networkError: "Error de red",
    serverError: "Error del servidor",
    unknownError: "Error desconocido",
    operationFailed: "Operación fallida",
    dataLoadFailed: "Error al cargar datos",
    saveFailed: "Error al guardar",
    deleteFailed: "Error al eliminar",
    copyFailed: "Error al copiar",
    invalidInput: "Entrada inválida",
    requiredField: "Este campo es obligatorio",
    loadChangeLogFailed: "Error al cargar el registro de cambios",
    loadDataFailed: "Error al cargar datos: ",
    selectExportPathFailed: "Error al seleccionar la ruta de exportación",
    exportFailed: "Error al exportar",
    unknownExportError: "Error desconocido durante la exportación",
  },

  // 更新日志
  changeLog: {
    new: "[Nuevo]",
    optimize: "[Optimización]",
    fix: "[Corrección]",
  },

  // 导出
  export: {
    verifyPasswordTitle: "Verificar contraseña de inicio de sesión",
    verifyPasswordDesc:
      "Por favor ingrese la contraseña de inicio de sesión de la bóveda actual para continuar con la operación de exportación.",
    loginPasswordLabel: "Contraseña de inicio de sesión",
    loginPasswordPlaceholder:
      "Por favor ingrese la contraseña de inicio de sesión",
    selectAccountsTitle: "Seleccionar cuentas a exportar",
    all: "Todo",
    byGroup: "Por grupo",
    byType: "Por tipo",
    manual: "Manual",
    selectGroup: "Seleccionar grupo",
    selectAll: "Seleccionar todo",
    clearAll: "Borrar todo",
    groupSelected: "Grupos seleccionados",
    accountsExpected: "cuentas esperadas",
    selectType: "Seleccionar tipo",
    typeSelected: "Tipos seleccionados",
    selectAccount: "Seleccionar cuenta",
    loadingAccountList: "Cargando lista de cuentas...",
    noAccountData: "No hay datos de cuenta",
    selectedAccountCount: "Cuentas seleccionadas",
    setBackupTitle: "Configurar contraseña de respaldo y ruta de exportación",
    backupPasswordLabel: "Contraseña de respaldo",
    backupPasswordPlaceholder: "Por favor ingrese la contraseña de respaldo",
    generate: "Generar",
    backupPasswordTip:
      "La contraseña de respaldo se usa para cifrar los datos exportados, guárdela de forma segura",
    exportPathLabel: "Ruta de exportación",
    exportPathPlaceholder: "Por favor seleccione la ruta de exportación",
    browse: "Navegar",
    exportingVault: "Exportando bóveda...",
    preparingExport: "Preparando exportación...",
    verifyingParameters: "Verificando parámetros...",
    exportComplete: "Exportación completada",
    exportSuccess: "Exportación exitosa",
    exportSuccessSubtitle: "La bóveda se ha exportado correctamente a: ",
    exportFailed: "Exportación fallida",
    openFolder: "Abrir carpeta",
    exportingWarning: "La exportación está en curso, espere...",
  },

  // Mensajes de éxito
  success: {
    operationSuccess: "Operación exitosa",
    dataSaved: "Datos guardados",
    dataDeleted: "Datos eliminados",
    copied: "Copiado",
    imported: "Importación exitosa",
    exported: "Exportación exitosa",
    dataRefreshed: "Datos actualizados",
    tabRenamed: "Pestaña renombrada exitosamente",
    tabDeleted: "Pestaña eliminada exitosamente",
    tabMoved: "Pestaña movida exitosamente",
    tabCreated: "Pestaña creada exitosamente",
    groupCreated: "Grupo creado exitosamente",
    groupRenamed: "Grupo renombrado exitosamente",
    groupDeleted: "Grupo eliminado exitosamente",
    groupMoved: "Grupo movido exitosamente",
    accountSaved: "Cuenta guardada exitosamente",
    accountDeleted: "Cuenta eliminada exitosamente",
    passwordGenerated: "Contraseña generada exitosamente",
    settingsSaved: "Configuración guardada exitosamente",
    vaultExported: "Bóveda exportada exitosamente",
    vaultImported: "Bóveda importada exitosamente",
    passwordVerified: "Contraseña verificada correctamente",
    backupPasswordGenerated: "Contraseña de respaldo generada",
    exportPathSelected: "Ruta de exportación seleccionada correctamente",
    vaultExportedSuccess: "La bóveda se exportó correctamente",
    usernameCopied:
      "Nombre de usuario copiado al portapapeles (se borrará automáticamente en 10 segundos)", // TODO: Translate
    passwordCopied: "Contraseña copiada al portapapeles", // TODO: Translate
    usernameAndPasswordCopied:
      "Nombre de usuario y contraseña copiados al portapapeles (se borrará automáticamente en 10 segundos)", // TODO: Translate
    urlCopied: "URL copiada al portapapeles", // TODO: Translate
    titleCopied: "Título copiado al portapapeles", // TODO: Translate
    notesCopied:
      "Notas copiadas al portapapeles (se borrará automáticamente en 10 segundos)", // TODO: Translate
  },

  // Advertencias
  warning: {
    noGroupData: "No hay datos de grupo, por favor cree un grupo primero",
    selectGroupFirst: "Por favor seleccione un grupo primero",
    selectTabFirst: "Por favor seleccione una pestaña primero",
    noAccountSelected: "Por favor seleccione una cuenta primero",
    confirmOperation: "Por favor confirme esta operación",
    selectAtLeastOneAccount: "Seleccione al menos una cuenta",
    exportInProgress: "La exportación está en curso, espere...",
    defaultGroupCannotBeDeleted: "El grupo predeterminado no se puede eliminar",
    groupNameCannotBeEmpty: "El nombre del grupo no puede estar vacío",
    groupNameNotChanged: "El nombre del grupo no ha cambiado",
    confirmDeleteGroup:
      '确定要删除分组"${groupName}"吗？\n删除后该分组下的所有标签和账号也将被删除，此操作不可恢复。',
    deleteGroupTitle: "Eliminar grupo",
    confirmDelete: "确定删除",
    cancelDelete: "取消",
    renameGroup: "Renombrar grupo",
    newGroupName: "Nuevo nombre de grupo",
    newGroupPlaceholder: "Nombre del grupo",
    createGroup: "Crear grupo",
    rename: "Renombrar",
    moveLeft: "Mover a la izquierda",
    moveRight: "Mover a la derecha",
  },

  // Estados
  status: {
    loading: "Cargando...",
    saving: "Guardando...",
    deleting: "Eliminando...",
    processing: "Procesando...",
    connecting: "Conectando...",
    searching: "Buscando...",
    exporting: "Exportando...",
    importing: "Importando...",
    generating: "Generando...",
    updatingAccountData: "Actualizando datos de la cuenta, espere...",
    totalAccounts: "Total de cuentas:",
    processed: "Procesado:",
    success: "Éxito:",
    failed: "Fallido:",
    startingChangeLoginPassword:
      "Iniciando cambio de contraseña de inicio de sesión...",
    loginPasswordChangeComplete:
      "Cambio de contraseña de inicio de sesión completado",
    changeFailed: "Cambio fallido",
  },

  // Menú más opciones
  moreMenu: {
    more: "Más",
    selectNewVault: "Seleccionar nueva bóveda",
    openVaultDirectory: "Abrir directorio de bóveda",
    generatePassword: "Generar contraseña",
    setPasswordRules: "Configurar reglas de generación de contraseña",
    changeLoginPassword: "Cambiar contraseña de inicio de sesión",
    exportVault: "Exportar bóveda",
    importVault: "Importar bóveda",
    changeLog: "Registro de cambios",
    settings: "Configuración",
    help: "Ayuda",
    about: "Acerca de",
    lockVault: "Bloquear bóveda",
    logout: "Cerrar sesión",
    oldLoginPasswordLabel: "Contraseña antigua",
    oldLoginPasswordPlaceholder:
      "Ingrese la contraseña de inicio de sesión actual",
    newLoginPasswordLabel: "Nueva contraseña",
    newLoginPasswordPlaceholder:
      "Ingrese la nueva contraseña de inicio de sesión",
    confirmNewPasswordLabel: "Confirmar nueva contraseña",
    confirmNewPasswordPlaceholder:
      "Vuelva a ingresar la nueva contraseña de inicio de sesión",
    selectVaultFile: "Seleccionar archivo de bóveda",
    selectVaultFilePrompt: "Por favor seleccione el archivo de bóveda a abrir:",
    selectVaultFilePlaceholder:
      "Por favor seleccione la ruta del archivo de bóveda",
    browse: "Navegar",
    supportedFileFormats: "Formatos de archivo compatibles: .db, .vault",
    openFileFailed: "Error al abrir el archivo",
    passwordCopied: "Contraseña copiada al portapapeles",
    copyFailedManual:
      "Error al copiar, por favor copie la contraseña manualmente",
    passwordChangeInProgress:
      "El cambio de contraseña está en curso, espere...",
    formRefNotFound: "No se encontró la referencia al formulario",
    oldPasswordIncorrect: "Contraseña antigua incorrecta, ingrese de nuevo",
    changePasswordConfirm:
      "Cambiar la contraseña de inicio de sesión volverá a cifrar todos los datos de la cuenta. Esta operación es irreversible, asegúrese de haber hecho una copia de seguridad de su bóveda actual. ¿Desea continuar?",
    passwordChangeSuccess:
      "¡Contraseña de inicio de sesión cambiada correctamente!",
    passwordChangeFailed: "Error al cambiar la contraseña de inicio de sesión",
    importSuccessDataUpdated: "Importación exitosa, datos actualizados",
    importSuccessDataRefreshFailed:
      "Importación exitosa, pero la actualización de datos falló, por favor actualice la página manualmente",
    vaultDirectoryOpened: "Directorio de la bóveda abierto",
    openDirectoryFailed: "Error al abrir el directorio",
    lockVaultConfirm:
      "¿Está seguro de que desea bloquear la bóveda? Después de bloquearla, deberá volver a ingresar la contraseña para acceder.",
    vaultLocked: "Bóveda bloqueada",
    lockVaultFailed: "Error al bloquear la bóveda",
    logoutConfirm:
      "¿Está seguro de que desea cerrar sesión? Después de cerrar sesión, deberá volver a ingresar la contraseña para acceder a la bóveda.",
    logoutSuccess: "Sesión cerrada",
    oldLoginPasswordRequired: "Ingrese la contraseña antigua",
    newLoginPasswordRequired: "Ingrese la nueva contraseña",
    newLoginPasswordMinLength: "La contraseña debe tener al menos 8 caracteres",
    newLoginPasswordStrength:
      "La contraseña debe contener mayúsculas, minúsculas, números y caracteres especiales",
    confirmNewLoginPasswordRequired: "Confirme la nueva contraseña",
    passwordsDoNotMatch: "Las contraseñas no coinciden",
    selectNewVaultConfirm:
      "Seleccionar una nueva bóveda lo devolverá a la pantalla de inicio de sesión, donde deberá volver a seleccionar el archivo de la bóveda e ingresar la contraseña. ¿Desea continuar?",
    selectNewVaultSuccess:
      "Ha vuelto a la pantalla de inicio de sesión, seleccione una nueva bóveda",
    pleaseSelectFile: "Por favor, seleccione un archivo primero",
    openingFile: "Abriendo archivo: ",
    open: "Abrir",
    continue: "Continuar",
    verifying: "Verificando...",
    changingPassword: "Cambiando...",
  },

  // Menú contextual
  contextMenu: {
    inputUsernameAndPassword: "Introducir nombre de usuario y contraseña",
    openUrl: "Abrir dirección",
    duplicate: "Generar copia",
    view: "Ver",
    edit: "Modificar",
    changeGroup: "Cambiar grupo",
    copyUsername: "Copiar cuenta",
    copyPassword: "Copiar contraseña",
    copyUsernameAndPassword: "Copiar cuenta y contraseña",
    copyUrl: "Copiar dirección",
    copyTitle: "Copiar título",
    copyNotes: "Copiar notas",
    delete: "Eliminar",
  },

  // Exportar
  export: {
    title: "Exportar bóveda",
    steps: {
      verifyPassword: "Verificar contraseña",
      selectAccounts: "Seleccionar cuentas",
      setBackup: "Configurar respaldo",
      exportComplete: "Exportación completa",
    },
    verifyPasswordTitle: "Verificar contraseña de inicio de sesión",
    verifyPasswordDesc:
      "Por favor ingrese la contraseña de inicio de sesión de la bóveda actual para continuar con la operación de exportación.",
    loginPassword: "Contraseña de inicio de sesión",
    loginPasswordPlaceholder:
      "Por favor ingrese la contraseña de inicio de sesión",
    selectAccountsTitle: "Seleccionar cuentas a exportar",
    exportAll: "Exportar todo",
    exportByGroup: "Exportar por grupo",
    exportByType: "Exportar por tipo", // TODO: Translate
    exportSelected: "Exportar seleccionados",
    selectGroups: "Seleccionar grupos", // TODO: Translate
    selectAll: "Seleccionar todo", // TODO: Translate
    clearAll: "Borrar todo", // TODO: Translate
    groupSelectionSummary:
      "Se han seleccionado {count} grupos, se espera exportar {accountCount} cuentas", // TODO: Translate
    selectTypes: "Seleccionar tipos", // TODO: Translate
    typeSelectionSummary:
      "Se han seleccionado {count} tipos, se espera exportar {accountCount} cuentas", // TODO: Translate
    selectAccounts: "Seleccionar cuentas", // TODO: Translate
    loadingAccounts: "Cargando cuentas...", // TODO: Translate
    noAccounts: "No hay datos de cuenta", // TODO: Translate
    setBackupTitle: "Establecer contraseña de respaldo y ruta de exportación", // TODO: Translate
    backupPassword: "Contraseña de respaldo", // TODO: Translate
    backupPasswordPlaceholder:
      "Por favor, introduzca la contraseña de respaldo", // TODO: Translate
    generate: "Generar", // TODO: Translate
    backupPasswordTip:
      "La contraseña de respaldo se utiliza para cifrar los datos exportados, por favor, guárdela de forma segura", // TODO: Translate
    exportPath: "Ruta de exportación", // TODO: Translate
    exportPathPlaceholder: "Por favor, seleccione la ruta de exportación", // TODO: Translate
    browse: "Navegar", // TODO: Translate
    exporting: "Exportando bóveda...", // TODO: Translate
    exportSuccessTitle: "Exportación exitosa", // TODO: Translate
    exportSuccessSubTitle: "La bóveda se ha exportado correctamente a: {path}", // TODO: Translate
    exportFailedTitle: "Error en la exportación", // TODO: Translate
    startExport: "Iniciar exportación", // TODO: Translate
    openFolder: "Abrir carpeta", // TODO: Translate
    loginPasswordRequired:
      "Por favor, introduzca la contraseña de inicio de sesión", // TODO: Translate
    backupPasswordRequired: "Por favor, introduzca la contraseña de respaldo", // TODO: Translate
    backupPasswordMinLength:
      "La contraseña de respaldo debe tener al menos 6 caracteres", // TODO: Translate
    exportPathRequired: "Por favor, seleccione la ruta de exportación", // TODO: Translate
  },

  // Importar
  import: {
    title: "Importar bóveda",
    selectFile: "Seleccionar archivo",
    selectFileDesc: "Por favor seleccione el archivo de bóveda a importar",
    fileFormat: "Formato de archivo",
    importProgress: "Progreso de importación",
    importComplete: "Importación completa",
  },

  // 新增的导入功能翻译
  importVault: {
    title: "Importar bóveda de contraseñas",
    step1: "Seleccionar archivo",
    step2: "Verificar contraseña",
    step3: "Importación completa",
    step1Title: "Seleccionar archivo de importación",
    step1Description:
      "Por favor seleccione el archivo de copia de seguridad de la bóveda de contraseñas (formato ZIP) para importar.",
    importFile: "Archivo de importación",
    selectImportFilePlaceholder:
      "Por favor seleccione el archivo de importación",
    browse: "Examinar",
    step2Title: "Introducir contraseña de descompresión",
    step2Description:
      "Por favor introduzca la contraseña de descompresión (contraseña de copia de seguridad) del archivo.",
    backupPassword: "Contraseña de descompresión",
    enterBackupPasswordPlaceholder:
      "Por favor introduzca la contraseña de descompresión",
    importingVault: "Importando bóveda de contraseñas...",
    importComplete: "Importación completa",
    importSuccess: "Importación exitosa",
    vaultImportedSuccessfully:
      "La bóveda de contraseñas se ha importado correctamente",
    importReport: "Informe de importación",
    totalAccounts: "Cuentas totales",
    successfullyImported: "Importadas con éxito",
    skippedAccounts: "Cuentas omitidas",
    errorAccounts: "Cuentas con error",
    totalGroups: "Grupos totales",
    importedGroups: "Grupos importados",
    totalTypes: "Tipos totales",
    importedTypes: "Tipos importados",
    skippedAccountDetails: "Detalles de cuentas omitidas",
    accountTitle: "Título de la cuenta",
    accountName: "Nombre de la cuenta",
    accountId: "ID de la cuenta",
    importFailed: "Importación fallida",
    importFailedStatus: "Importación fallida",
    unknownError: "Se produjo un error desconocido durante la importación",
    close: "Cerrar",
    cancel: "Cancelar",
    previousStep: "Paso anterior",
    startImport: "Iniciar importación",
    nextStep: "Siguiente paso",
    refreshData: "Actualizar datos",
    selectImportFileMessage: "Por favor seleccione el archivo de importación",
    enterBackupPasswordMessage:
      "Por favor introduzca la contraseña de descompresión",
    selectFileSuccess: "Archivo de importación seleccionado correctamente",
    selectFileFailed: "Error al seleccionar el archivo de importación",
    preparingToImport: "Preparando para importar...",
    validatingFileAndPassword: "Validando archivo y contraseña...",
    vaultImportSuccess: "Bóveda de contraseñas importada correctamente",
    dataRefreshed: "Datos actualizados",
    importInProgressWarning: "La importación está en curso, espere...",
  },

  // 密码生成器
  passwordGenerator: {
    title: "Generador de contraseñas",
    selectRule: "Seleccionar regla de contraseña",
    selectRulePlaceholder: "Por favor seleccione una regla de contraseña",
    generalRule: "Regla de contraseña general",
    customRule: "Regla de contraseña personalizada",
    includeUppercase: "Incluir mayúsculas",
    includeLowercase: "Incluir minúsculas",
    includeNumbers: "Incluir números",
    includeSpecialChars: "Incluir caracteres especiales",
    passwordLength: "Longitud de contraseña",
    customSpecialChars: "Caracteres especiales personalizados",
    defaultSpecialChars:
      "Dejar vacío para usar caracteres especiales por defecto",
    generatePassword: "Generar contraseña",
    usePassword: "Usar esta contraseña",
    generatedPassword: "Contraseña generada",
    clickToGenerate: "Haga clic en el botón generar contraseña",
    selectCharType: "Por favor seleccione al menos un tipo de carácter",
    enterPattern: "Por favor ingrese el patrón de contraseña",
    generateFirst: "Por favor genere una contraseña primero",
  },

  // 新增的密码规则设置翻译
  passwordRuleSettings: {
    title: "Configurar reglas de generación de contraseña",
    savedRules: "Reglas de contraseña guardadas",
    newRule: "Nueva regla",
    ruleName: "Nombre de la regla",
    description: "Descripción",
    operation: "Operación",
    edit: "Editar",
    delete: "Eliminar",
    generalRule: "Regla de contraseña general",
    includeUppercase: "Incluir mayúsculas",
    includeLowercase: "Incluir minúsculas",
    includeNumbers: "Incluir números",
    includeSpecialChars: "Incluir caracteres especiales",
    passwordLength: "Longitud de contraseña",
    customSpecialChars: "Caracteres especiales personalizados",
    customSpecialCharsPlaceholder:
      "Dejar en blanco para usar caracteres especiales predeterminados: !@#$%^&*()_+-=[]{}|;:,.<>?",
    customRulesDescription:
      "Descripción de reglas de contraseña personalizadas",
    close: "Cerrar",
    saveSettings: "Guardar configuración",
    editRule: "Editar regla",
    newRuleDialog: "Nueva regla",
    ruleNameRequired: "Por favor introduzca el nombre de la regla",
    ruleDescriptionPlaceholder:
      "Por favor introduzca la descripción de la regla",
    ruleType: "Tipo de regla",
    general: "Regla general",
    custom: "Regla personalizada",
    passwordPattern: "Patrón de contraseña",
    passwordPatternPlaceholder: "Ejemplo: Aaa111",
    save: "Guardar",
    cancel: "Cancelar",
    confirmDeleteRule:
      '¿Está seguro de que desea eliminar la regla "{ruleName}"?',
    deleteConfirmation: "Confirmación de eliminación",
    confirmDelete: "Confirmar",
    ruleDeletedSuccess: "Regla eliminada correctamente",
    ruleUpdateSuccess: "Regla actualizada correctamente",
    ruleCreateSuccess: "Regla creada correctamente",
    enterRuleName: "Por favor introduzca el nombre de la regla",
    settingsSaved: "Configuración guardada correctamente",
  },

  // 新增的搜索结果翻译
  searchResults: {
    groups: "Grupos",
    accounts: "Cuentas",
    noUrl: "Sin URL",
    belongsToGroup: "Pertenece al grupo",
    unknownGroup: "Grupo desconocido",
    noSearchResults: "No se encontraron resultados de búsqueda",
    tryOtherKeywords: "Intente con otras palabras clave",
  },

  // 新增的状态栏翻译
  statusBar: {
    authInfo: "Información de autorización: Versión Pro activada",
    usageCount: "Usos: {count} veces",
    usageDays: "Días de uso: {days} días",
  },

  // 新增的标签上下文菜单翻译
  tabContextMenu: {
    rename: "Renombrar",
    deleteTab: "Eliminar pestaña",
    newTab: "Nueva pestaña",
    moveUp: "Mover arriba",
    moveDown: "Mover abajo",
    promptNewName: "Por favor introduzca el nuevo nombre de la pestaña",
    renameTab: "Renombrar pestaña",
    confirm: "Confirmar",
    cancel: "Cancelar",
    tabName: "Nombre de la pestaña",
    tabNameCannotBeEmpty: "El nombre de la pestaña no puede estar vacío",
    tabNameNotChanged: "El nombre de la pestaña no ha cambiado",
    userCanceledRename: "Usuario canceló la operación de renombrado",
    confirmDeleteTab:
      '¿Está seguro de que desea eliminar la pestaña "{tabName}"?\nDespués de eliminarla, todas las cuentas bajo esta pestaña también se eliminarán, esta operación es irreversible.',
    deleteConfirmation: "Eliminar pestaña",
    confirmDelete: "Confirmar eliminación",
    userCanceledDelete: "Usuario canceló la operación de eliminación",
    promptNewTabName: "Por favor introduzca el nuevo nombre de la pestaña",
    userCanceledNewTab: "Usuario canceló la operación de nueva pestaña",
  },

  // 新增的标签侧边栏翻译
  tabsSidebar: {
    emptyTabs: "No hay pestañas",
    createTabHint: "Haga clic en el botón de abajo para crear una pestaña",
    newTab: "Nueva pestaña",
  },

  // 新增的测试对话框翻译
  testDialog: {
    title: "Diálogo de prueba",
    description1:
      "Este es un diálogo de prueba para verificar si la función de ventana emergente funciona correctamente.",
    description2:
      "Si puede ver este diálogo, significa que la función de ventana emergente funciona correctamente.",
    close: "Cerrar",
  },

  // 新增的标题栏翻译
  titleBar: {
    defaultAppTitle: "Gestor de contraseñas",
  },

  // 新增的锁定事件服务翻译
  lockEventService: {
    logPrefix: "Servicio de eventos de bloqueo",
    frontendStateUpdated: "Estado frontend actualizado",
    redirectToLogin: "Redirigir a la página de inicio de sesión",
    vaultAutoLocked:
      "La bóveda de contraseñas se ha bloqueado automáticamente, por favor inicie sesión de nuevo",
    sensitiveDataCleared: "Datos sensibles eliminados",
    clearSensitiveDataFailed: "Error al eliminar datos sensibles",
    windowEventListenersSet:
      "Se han configurado los oyentes de eventos de ventana",
    windowLostFocus: "Ventana perdió el foco",
    windowGainedFocus: "Ventana ganó el foco",
    backendNotifiedMinimize:
      "Se notificó al backend la minimización de la ventana",
    notifyMinimizeFailed: "Error al notificar la minimización de la ventana",
    backendNotifiedFocus: "Se notificó al backend que la ventana ganó el foco",
    notifyFocusFailed: "Error al notificar que la ventana ganó el foco",
    manualLockCheckTriggered: "Comprobación de bloqueo manual activada",
  },

  // 新增的账号工具函数翻译
  accountUtils: {
    defaultAccountTitle: "Cuenta",
    startCopyingPassword:
      "Comenzando a copiar contraseña, ID de cuenta: {accountId}",
    passwordCopiedSuccess:
      "Contraseña copiada al portapapeles (se borrará automáticamente en 10 segundos)",
    passwordCopySuccess:
      "Contraseña copiada con éxito, ID de cuenta: {accountId}",
    passwordCopyFailed: "Error al copiar contraseña",
    startGettingPassword:
      "Comenzando a obtener contraseña, ID de cuenta: {accountId}",
    passwordGetSuccess:
      "Contraseña obtenida con éxito, ID de cuenta: {accountId}",
    passwordGetFailed: "Error al obtener contraseña",
    startGettingAccountDetail:
      "Comenzando a obtener detalles de la cuenta, ID de cuenta: {accountId}",
    accountDetailGetSuccess:
      "Detalles de la cuenta obtenidos con éxito, ID de cuenta: {accountId}",
    accountDetailGetFailed: "Error al obtener detalles de la cuenta",
  },
};
