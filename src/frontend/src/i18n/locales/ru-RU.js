/**
 * Русский языковой пакет
 * @author Чэнь Фэнцин
 * @date 2025-10-05
 */

export default {
  // Общие
  common: {
    confirm: "Подтвердить",
    cancel: "Отмена",
    save: "Сохранить",
    delete: "Удалить",
    edit: "Редактировать",
    add: "Добавить",
    search: "Поиск",
    close: "Закрыть",
    ok: "ОК",
    yes: "Да",
    no: "Нет",
    loading: "Загрузка...",
    success: "Успех",
    error: "Ошибка",
    warning: "Предупреждение",
    info: "Информация",
    copy: "Копировать",
    paste: "Вставить",
    cut: "Вырезать",
    refresh: "Обновить",
    reset: "Сбросить",
    clear: "Очистить",
    back: "Назад",
    next: "Далее",
    previous: "Предыдущий",
    finish: "Завершить",
    submit: "Отправить",
    import: "Импорт",
    export: "Экспорт",
    settings: "Настройки",
    help: "Помощь",
    about: "О программе",
    more: "Больше",
  },

  // Страница входа
  login: {
    title: "Вход в хранилище",
    password: "Пароль",
    passwordPlaceholder: "Введите пароль хранилища",
    login: "Войти",
    createVault: "Создать хранилище",
    openVault: "Открыть другое хранилище",
    invalidPassword: "Неверный пароль",
    loginSuccess: "Успешный вход",
    loginFailed: "Ошибка входа",
    recentVaults: "Недавно использованные",
    passwordHint1:
      "Пароль для входа является ключевым для шифрования и дешифрования хранилища",
    passwordHint2:
      "Пожалуйста, храните свой пароль в безопасности, после потери его невозможно восстановить",
    passwordHint3:
      "Пароль для входа хранится в зашифрованном виде в базе данных",
    vaultFile: "Файл хранилища",
    vaultFilePlaceholder: "Пожалуйста, выберите файл хранилища",
    vaultPathRequired: "Пожалуйста, введите путь к файлу хранилища",
    passwordRequired: "Пожалуйста, введите пароль для входа",
    passwordMinLength: "Длина пароля не может быть меньше 6 символов",
    selectFileFailed: "Не удалось выбрать файл",
    passwordEmpty: "Пароль не может быть пустым",
    passwordMinLength8: "Длина пароля должна быть не менее 8 символов",
    passwordNeedUppercase: "Пароль должен содержать заглавные буквы",
    passwordNeedLowercase: "Пароль должен содержать строчные буквы",
    passwordNeedNumber: "Пароль должен содержать цифры",
    passwordNeedSpecialChar:
      "Пароль должен содержать специальные символы (!@#$%^&* и т. д.)",
    passwordStrengthOk: "Надежность пароля соответствует требованиям",
    vaultFileNotExists: "Файл хранилища не существует",
    createVaultPrompt:
      "Пожалуйста, введите имя хранилища (без расширения)\n\nСистема автоматически создаст каталог данных в текущем каталоге программы и сохранит файл хранилища в каталоге данных, автоматически создавая файл .db в соответствии с именем.",
    vaultNamePlaceholder: "Пример: мое_хранилище",
    setPasswordPrompt:
      "Пожалуйста, установите пароль для входа\n\nТребования к паролю:\n• Не менее 8 символов\n• Содержит заглавные буквы\n• Содержит строчные буквы\n• Содержит цифры\n• Содержит специальные символы (!@#$%^&* и т. д.)",
    setPassword: "Установить пароль",
    passwordRequirements: "Требования к паролю",
    passwordReq1: "Не менее 8 символов",
    passwordReq2: "Содержит заглавные буквы",
    passwordReq3: "Содержит строчные буквы",
    passwordReq4: "Содержит цифры",
    passwordReq5: "Содержит специальные символы (!@#$%^&* и т. д.)",
    createVaultSuccess: "Хранилище успешно создано",
    createVaultFailed: "Не удалось создать хранилище",
    switchToFullMode:
      "Переключено в полный режим, пожалуйста, выберите файл хранилища",
    switchToCreateMode:
      "Переключено в полный режим, пожалуйста, создайте новое хранилище",
  },

  // Главный интерфейс
  main: {
    title: "wepass - Менеджер паролей",
    searchPlaceholder: "Поиск аккаунтов...",
    noData: "Нет данных",
    noSearchResults: "Соответствующие результаты не найдены",
    searchResults: "Результаты поиска",
    groups: "Группы",
    group: "Группа",
    accounts: "Аккаунты",
    allAccounts: "Все аккаунты",
    noAddress: "Нет адреса",
    belongsToGroup: "Принадлежит группе",
    inputUsernameAndPassword: "Ввести имя пользователя и пароль",
    inputUsername: "Ввести имя пользователя",
    inputPassword: "Ввести пароль",
    tryOtherKeywords: "Попробуйте другие ключевые слова",
    unknownGroup: "Неизвестная группа",
  },

  // Диалог настроек
  settings: {
    title: "Настройки",
    general: "Общие",
    log: "Журнал",
    lock: "Блокировка",
    theme: "Тема",
    language: "Язык",
    lightTheme: "Светлая",
    darkTheme: "Тёмная",
    selectTheme: "Выберите тему",
    selectLanguage: "Выберите язык",
    logSettings: "Настройки журнала",
    enableInfoLog: "Включить информационный журнал",
    enableDebugLog: "Включить журнал отладки",
    lockSettings: "Настройки блокировки",
    enableAutoLock: "Включить автоблокировку",
    enableTimerLock: "Включить блокировку по таймеру",
    enableMinimizeLock: "Включить блокировку при сворачивании",
    lockTimeMinutes: "Время блокировки (минуты)",
    enableSystemLock: "Включить системную блокировку",
    systemLockMinutes: "Время системной блокировки (минуты)",
    settingsSaved: "Настройки сохранены",
    loadSettingsFailed: "Не удалось загрузить настройки",
    saveSettingsFailed: "Не удалось сохранить настройки",
  },

  // Управление аккаунтами
  account: {
    title: "Название",
    username: "Имя пользователя",
    password: "Пароль",
    url: "URL",
    type: "Тип",
    notes: "Заметки",
    icon: "Иконка",
    favorite: "Избранное",
    useCount: "Количество использований",
    lastUsed: "Последнее использование",
    created: "Создано",
    updated: "Обновлено",
    inputMethod: "Метод ввода",
    addAccount: "Добавить аккаунт",
    editAccount: "Редактировать аккаунт",
    deleteAccount: "Удалить аккаунт",
    copyUsername: "Копировать имя пользователя",
    copyPassword: "Копировать пароль",
    openUrl: "Открыть URL",
    generatePassword: "Сгенерировать пароль",
    showPassword: "Показать пароль",
    hidePassword: "Скрыть пароль",
    accountSaved: "Аккаунт сохранён",
    accountDeleted: "Аккаунт удалён",
    confirmDelete: "Вы уверены, что хотите удалить этот аккаунт?",
    deleteConfirmTitle: "Подтверждение удаления",
    copiedToClipboard: "Скопировано в буфер обмена",
    copyFailed: "Ошибка копирования",
    noAccounts: "Нет аккаунтов",
    noAccountsHint: "Нажмите кнопку '+' для создания первого аккаунта",
  },

  // Группы и типы
  group: {
    allGroups: "Все группы",
    addGroup: "Добавить группу",
    editGroup: "Редактировать группу",
    deleteGroup: "Удалить группу",
    groupName: "Название группы",
    groupIcon: "Иконка группы",
    rename: "Переименовать", // TODO: Translate
    createGroup: "Создать группу", // TODO: Translate
    moveLeft: "Переместить влево", // TODO: Translate
    moveRight: "Переместить вправо", // TODO: Translate
    newGroupName: "Новое имя группы", // TODO: Translate
    renameGroup: "Переименовать группу", // TODO: Translate
    newGroupPlaceholder: "Имя группы", // TODO: Translate
    groupNameCannotBeEmpty: "Имя группы не может быть пустым", // TODO: Translate
    groupNameNotChanged: "Имя группы не изменилось", // TODO: Translate
    defaultGroupCannotBeDeleted: "Группу по умолчанию нельзя удалить", // TODO: Translate
    confirmDeleteGroup:
      'Вы уверены, что хотите удалить группу "{groupName}"?\nВсе теги и учетные записи в этой группе также будут удалены, это действие необратимо.', // TODO: Translate
    deleteGroupTitle: "Удалить группу", // TODO: Translate
    confirmDelete: "Подтвердить удаление", // TODO: Translate
    defaultGroupName: "По умолчанию", // TODO: Translate
  },

  type: {
    allTypes: "Все типы",
    addType: "Добавить тип",
    editType: "Редактировать тип",
    deleteType: "Удалить тип",
    typeName: "Название типа",
    typeIcon: "Иконка типа",
  },

  // Диалог помощи
  help: {
    title: "Помощь",
    quickStart: "Быстрый старт",
    welcome: "Добро пожаловать в WePassword",
    welcomeTitle: "Добро пожаловать в WePassword",
    description:
      "WePassword - это безопасный инструмент управления паролями, который поможет вам безопасно хранить и управлять паролями.",
    welcomeDescription:
      "WePassword - это безопасный инструмент управления паролями, который поможет вам безопасно хранить и управлять паролями.",
    basicOperations: "Основные операции",
    createVault:
      "Создать хранилище: При первом использовании система поможет вам создать новое хранилище",
    addPassword:
      'Добавить пароль: Нажмите кнопку "Добавить", чтобы ввести информацию о сайте, имени пользователя и пароле',
    searchPassword:
      "Поиск пароля: Используйте поле поиска для быстрого поиска конкретных записей паролей",
    editPassword:
      "Редактировать пароль: Дважды щёлкните запись пароля или нажмите кнопку редактирования для изменения",
    deletePassword:
      "Удалить пароль: Выберите запись пароля и нажмите кнопку удаления",
    features: "Функции", // TODO: Translate
    mainFeatures: "Основные функции", // TODO: Translate
    passwordGenerator: "Генератор паролей", // TODO: Translate
    passwordGeneratorDesc:
      "Генерируйте надежные пароли, поддерживайте пользовательские правила", // TODO: Translate
    groupManagement: "Управление группами", // TODO: Translate
    groupManagementDesc:
      "Используйте вкладки для классификации и управления паролями", // TODO: Translate
    secureEncryption: "Безопасное шифрование", // TODO: Translate
    secureEncryptionDesc: "Все данные паролей хранятся в зашифрованном виде", // TODO: Translate
    importExport: "Импорт/Экспорт", // TODO: Translate
    importExportDesc: "Поддержка импорта данных из других менеджеров паролей", // TODO: Translate
    backupRestore: "Резервное копирование/Восстановление", // TODO: Translate
    backupRestoreDesc:
      "Регулярно создавайте резервные копии своего хранилища паролей", // TODO: Translate
    passwordGenerationRules: "Правила генерации паролей", // TODO: Translate
    passwordGenerationRulesDesc:
      "Поддержка нескольких наборов символов и пользовательских правил:", // TODO: Translate
    lowercaseLetters: "Строчные буквы", // TODO: Translate
    mixedCaseLetters: "Буквы в смешанном регистре", // TODO: Translate
    uppercaseLetters: "Заглавные буквы", // TODO: Translate
    digits: "Цифры", // TODO: Translate
    specialCharacters: "Специальные символы", // TODO: Translate
    customCharacterSet: "Пользовательский набор символов", // TODO: Translate
    securityTips: "Советы по безопасности", // TODO: Translate
    securitySuggestions: "Рекомендации по безопасности", // TODO: Translate
    masterPassword: "Мастер-пароль", // TODO: Translate
    masterPasswordDesc: "Установите надежный мастер-пароль и запомните его", // TODO: Translate
    regularBackup: "Регулярное резервное копирование", // TODO: Translate
    regularBackupDesc:
      "Регулярно создавайте резервные копии файла вашего хранилища паролей", // TODO: Translate
    timelyUpdate: "Своевременное обновление", // TODO: Translate
    timelyUpdateDesc: "Регулярно обновляйте пароли для важных учетных записей", // TODO: Translate
    avoidRepetition: "Избегайте повторений", // TODO: Translate
    avoidRepetitionDesc:
      "Не используйте один и тот же пароль на нескольких веб-сайтах", // TODO: Translate
    safeEnvironment: "Безопасная среда", // TODO: Translate
    safeEnvironmentDesc: "Используйте менеджер паролей в безопасной среде", // TODO: Translate
    precautions: "Меры предосторожности", // TODO: Translate
    precaution1:
      "Пожалуйста, храните свой мастер-пароль в безопасности. Если вы его забудете, вы не сможете восстановить свои данные.", // TODO: Translate
    precaution2:
      "Рекомендуется периодически менять мастер-пароль для повышения безопасности.", // TODO: Translate
    precaution3:
      "Не сохраняйте файл хранилища паролей на общедоступном компьютере.", // TODO: Translate
    faq: "Часто задаваемые вопросы", // TODO: Translate
    faqTitle: "Часто задаваемые вопросы", // TODO: Translate
    faq1_q: "Что делать, если я забуду свой мастер-пароль?", // TODO: Translate
    faq1_a:
      "К сожалению, если вы забудете свой мастер-пароль, вы не сможете восстановить данные в хранилище паролей. Рекомендуется регулярно создавать резервные копии своего хранилища паролей и записывать свой мастер-пароль в безопасном месте.", // TODO: Translate
    faq2_q: "Как создать резервную копию хранилища паролей?", // TODO: Translate
    faq2_a:
      'Вы можете скопировать файл хранилища паролей в безопасное место или использовать функцию экспорта в меню "Больше".', // TODO: Translate
    faq3_q: "Какие форматы файлов поддерживаются?", // TODO: Translate
    faq3_a:
      "В настоящее время поддерживаются файлы хранилища паролей в формате .db и .vault.", // TODO: Translate
    faq4_q: "Как сгенерировать надежный пароль?", // TODO: Translate
    faq4_a:
      'Используйте функцию "Сгенерировать пароль" в меню "Больше". Рекомендуется использовать комбинацию прописных и строчных букв, цифр и специальных символов длиной не менее 12 символов.', // TODO: Translate
  },

  // Диалог "О программе"
  about: {
    title: "О программе",
    appName: "Менеджер паролей",
    description:
      "Кроссплатформенный инструмент управления паролями на основе Wails + Go + Vue.js",
    version: "Версия",
    buildDate: "Дата сборки",
    author: "Автор",
    support: "Поддержка",
    license: "Лицензия",
    github: "GitHub",
    gitee: "Gitee",
  },

  // Сообщения об ошибках
  error: {
    networkError: "Ошибка сети",
    serverError: "Ошибка сервера",
    unknownError: "Неизвестная ошибка",
    operationFailed: "Операция не удалась",
    dataLoadFailed: "Не удалось загрузить данные",
    saveFailed: "Не удалось сохранить",
    deleteFailed: "Не удалось удалить",
    copyFailed: "Не удалось скопировать",
    invalidInput: "Неверный ввод",
    requiredField: "Это поле обязательно для заполнения",
    apiServiceUnavailable: "API-сервис недоступен",
    backendServiceUnavailable:
      "Бэкенд-сервис недоступен, пожалуйста, проверьте, работает ли приложение правильно",
    accountDataFormatError: "Ошибка формата данных аккаунта",
    accessibilityPermissionRequired:
      "Требуется разрешение на доступность для автозаполнения",
    autofillFailed: "Автозаполнение не удалось",
    autofillUsernameFailed: "Не удалось автозаполнить имя пользователя",
    autofillPasswordFailed: "Не удалось автозаполнить пароль",
  },

  // Сообщения об успехе
  success: {
    operationSuccess: "Операция выполнена успешно",
    dataSaved: "Данные сохранены",
    dataDeleted: "Данные удалены",
    copied: "Скопировано",
    imported: "Импорт выполнен успешно",
    exported: "Экспорт выполнен успешно",
    dataRefreshed: "Данные обновлены",
    tabRenamed: "Вкладка успешно переименована",
    tabDeleted: "Вкладка успешно удалена",
    tabMoved: "Вкладка успешно перемещена",
    tabCreated: "Вкладка успешно создана",
    groupCreated: "Группа успешно создана",
    groupRenamed: "Группа успешно переименована",
    groupDeleted: "Группа успешно удалена",
    groupMoved: "Группа успешно перемещена",
    accountSaved: "Аккаунт успешно сохранён",
    accountDeleted: "Аккаунт успешно удалён",
    passwordGenerated: "Пароль успешно сгенерирован",
    settingsSaved: "Настройки успешно сохранены",
    vaultExported: "Хранилище успешно экспортировано",
    vaultImported: "Хранилище успешно импортировано",
    autofillUsernameAndPassword:
      "Имя пользователя и пароль успешно автозаполнены",
    autofillUsername: "Имя пользователя успешно автозаполнено",
    autofillPassword: "Пароль успешно автозаполнен",
  },

  // Предупреждающие сообщения
  warning: {
    noGroupData: "Нет данных группы, сначала создайте группу",
    selectGroupFirst: "Сначала выберите группу",
    selectTabFirst: "Сначала выберите вкладку",
    noAccountSelected: "Сначала выберите аккаунт",
    confirmOperation: "Пожалуйста, подтвердите эту операцию",
  },

  // Информация о статусе
  status: {
    loading: "Загрузка...",
    saving: "Сохранение...",
    deleting: "Удаление...",
    processing: "Обработка...",
    connecting: "Подключение...",
    searching: "Поиск...",
    exporting: "Экспорт...",
    importing: "Импорт...",
    generating: "Генерация...",
  },

  // Меню "Больше"
  moreMenu: {
    more: "Больше",
    selectNewVault: "Выбрать новое хранилище",
    openVaultDirectory: "Открыть каталог хранилища",
    generatePassword: "Сгенерировать пароль",
    setPasswordRules: "Установить правила генерации пароля",
    changeLoginPassword: "Изменить пароль входа",
    exportVault: "Экспортировать хранилище",
    importVault: "Импортировать хранилище",
    changeLog: "Журнал изменений",
    settings: "Настройки",
    help: "Помощь",
    about: "О программе",
    lockVault: "Заблокировать хранилище",
    logout: "Выйти",
    oldLoginPasswordLabel: "Старый пароль",
    oldLoginPasswordPlaceholder: "Введите текущий пароль для входа",
    newLoginPasswordLabel: "Новый пароль",
    newLoginPasswordPlaceholder: "Введите новый пароль для входа",
    confirmNewPasswordLabel: "Подтвердить новый пароль",
    confirmNewPasswordPlaceholder: "Повторно введите новый пароль для входа",
    selectVaultFile: "Выбрать файл хранилища",
    selectVaultFilePrompt: "Пожалуйста, выберите файл хранилища для открытия:",
    selectVaultFilePlaceholder: "Пожалуйста, выберите путь к файлу хранилища",
    browse: "Обзор",
    supportedFileFormats: "Поддерживаемые форматы файлов: .db, .vault",
    openFileFailed: "Не удалось открыть файл",
    passwordCopied: "Пароль скопирован в буфер обмена",
    copyFailedManual:
      "Не удалось скопировать, пожалуйста, скопируйте пароль вручную",
    passwordChangeInProgress:
      "Смена пароля в процессе, пожалуйста, подождите...",
    formRefNotFound: "Ссылка на форму не найдена",
    oldPasswordIncorrect: "Старый пароль неверен, введите снова",
    changePasswordConfirm:
      "Смена пароля для входа повторно зашифрует все данные аккаунта. Эта операция необратима, убедитесь, что вы сделали резервную копию вашего текущего хранилища. Вы хотите продолжить?",
    passwordChangeSuccess: "Пароль для входа успешно изменен!",
    passwordChangeFailed: "Не удалось изменить пароль для входа",
    importSuccessDataUpdated: "Импорт успешно завершен, данные обновлены",
    importSuccessDataRefreshFailed:
      "Импорт успешно завершен, но обновление данных не удалось, пожалуйста, обновите страницу вручную",
    vaultDirectoryOpened: "Каталог хранилища открыт",
    openDirectoryFailed: "Не удалось открыть каталог",
    lockVaultConfirm:
      "Вы уверены, что хотите заблокировать хранилище? После блокировки вам нужно будет снова ввести пароль для доступа.",
    vaultLocked: "Хранилище заблокировано",
    lockVaultFailed: "Не удалось заблокировать хранилище",
    logoutConfirm:
      "Вы уверены, что хотите выйти из системы? После выхода вам нужно будет снова ввести пароль для доступа к хранилищу.",
    logoutSuccess: "Выход выполнен успешно",
    oldLoginPasswordRequired: "Введите старый пароль",
    newLoginPasswordRequired: "Введите новый пароль",
    newLoginPasswordMinLength: "Пароль должен содержать не менее 8 символов",
    newLoginPasswordStrength:
      "Пароль должен содержать заглавные буквы, строчные буквы, цифры и специальные символы",
    confirmNewLoginPasswordRequired: "Подтвердите новый пароль",
    passwordsDoNotMatch: "Пароли не совпадают",
    selectNewVaultConfirm:
      "Выбор нового хранилища вернет вас на экран входа, где вам нужно будет снова выбрать файл хранилища и ввести пароль. Вы хотите продолжить?",
    selectNewVaultSuccess:
      "Вы вернулись на экран входа, выберите новое хранилище",
    pleaseSelectFile: "Пожалуйста, сначала выберите файл",
    openingFile: "Открытие файла: ",
    open: "Открыть",
    continue: "Продолжить",
    verifying: "Проверка...",
    changingPassword: "Изменение...",
  },

  // Контекстное меню
  contextMenu: {
    inputUsernameAndPassword: "Ввести имя пользователя и пароль", // TODO: Translate
    openUrl: "Открыть адрес",
    duplicate: "Создать копию",
    view: "Просмотр",
    edit: "Редактировать",
    changeGroup: "Изменить группу",
    copyUsername: "Копировать имя пользователя",
    copyPassword: "Копировать пароль",
    copyUsernameAndPassword: "Копировать имя пользователя и пароль",
    copyUrl: "Копировать адрес",
    copyTitle: "Копировать заголовок",
    copyNotes: "Копировать заметки",
    delete: "Удалить",
  },

  // Функция экспорта
  export: {
    title: "Экспорт хранилища",
    steps: {
      verifyPassword: "Проверить пароль",
      selectAccounts: "Выбрать аккаунты",
      setBackup: "Настроить резервную копию",
      exportComplete: "Экспорт завершён",
    },
    verifyPasswordTitle: "Проверить пароль входа",
    verifyPasswordDesc:
      "Пожалуйста, введите текущий пароль входа в хранилище для продолжения операции экспорта.",
    loginPassword: "Пароль входа",
    loginPasswordPlaceholder: "Пожалуйста, введите пароль входа",
    selectAccountsTitle: "Выберите аккаунты для экспорта",
    exportAll: "Экспортировать всё",
    exportByGroup: "Экспорт по группам",
    exportByType: "Экспорт по типу", // TODO: Translate
    exportSelected: "Экспорт выбранного",
    selectGroups: "Выбрать группы", // TODO: Translate
    selectAll: "Выбрать все", // TODO: Translate
    clearAll: "Снять выделение", // TODO: Translate
    groupSelectionSummary:
      "Выбрано {count} групп, ожидается экспорт {accountCount} аккаунтов", // TODO: Translate
    selectTypes: "Выбрать типы", // TODO: Translate
    typeSelectionSummary:
      "Выбрано {count} типов, ожидается экспорт {accountCount} аккаунтов", // TODO: Translate
    selectAccounts: "Выбрать аккаунты", // TODO: Translate
    loadingAccounts: "Загрузка аккаунтов...", // TODO: Translate
    noAccounts: "Нет данных аккаунта", // TODO: Translate
    setBackupTitle: "Установить пароль для резервной копии и путь экспорта", // TODO: Translate
    backupPassword: "Пароль резервной копии", // TODO: Translate
    backupPasswordPlaceholder: "Пожалуйста, введите пароль резервной копии", // TODO: Translate
    generate: "Сгенерировать", // TODO: Translate
    backupPasswordTip:
      "Пароль резервной копии используется для шифрования экспортированных данных, пожалуйста, храните его в безопасности", // TODO: Translate
    exportPath: "Путь экспорта", // TODO: Translate
    exportPathPlaceholder: "Пожалуйста, выберите путь экспорта", // TODO: Translate
    browse: "Обзор", // TODO: Translate
    exporting: "Экспорт хранилища...", // TODO: Translate
    exportSuccessTitle: "Экспорт успешен", // TODO: Translate
    exportSuccessSubTitle: "Хранилище успешно экспортировано в: {path}", // TODO: Translate
    exportFailedTitle: "Экспорт не удался", // TODO: Translate
    startExport: "Начать экспорт", // TODO: Translate
    openFolder: "Открыть папку", // TODO: Translate
    loginPasswordRequired: "Пожалуйста, введите пароль для входа", // TODO: Translate
    backupPasswordRequired: "Пожалуйста, введите пароль резервной копии", // TODO: Translate
    backupPasswordMinLength:
      "Пароль резервной копии должен содержать не менее 6 символов", // TODO: Translate
    exportPathRequired: "Пожалуйста, выберите путь экспорта", // TODO: Translate
  },

  // Функция импорта
  import: {
    title: "Импорт хранилища",
    selectFile: "Выбрать файл",
    selectFileDesc: "Пожалуйста, выберите файл хранилища для импорта",
    fileFormat: "Формат файла",
    importProgress: "Прогресс импорта",
    importComplete: "Импорт завершён",
  },

  // Новые переводы для функции импорта
  importVault: {
    title: "Импорт хранилища паролей",
    step1: "Выбрать файл",
    step2: "Проверить пароль",
    step3: "Импорт завершён",
    step1Title: "Выбрать файл для импорта",
    step1Description:
      "Пожалуйста, выберите файл резервной копии хранилища паролей (формат ZIP) для импорта.",
    importFile: "Файл для импорта",
    selectImportFilePlaceholder: "Пожалуйста, выберите файл для импорта",
    browse: "Обзор",
    step2Title: "Введите пароль для распаковки",
    step2Description:
      "Пожалуйста, введите пароль для распаковки (резервный пароль) файла.",
    backupPassword: "Пароль для распаковки",
    enterBackupPasswordPlaceholder: "Пожалуйста, введите пароль для распаковки",
    importingVault: "Импорт хранилища паролей...",
    importComplete: "Импорт завершён",
    importSuccess: "Импорт успешно завершён",
    vaultImportedSuccessfully: "Хранилище паролей успешно импортировано",
    importReport: "Отчёт об импорте",
    totalAccounts: "Всего аккаунтов",
    successfullyImported: "Успешно импортировано",
    skippedAccounts: "Пропущенные аккаунты",
    errorAccounts: "Аккаунты с ошибками",
    totalGroups: "Всего групп",
    importedGroups: "Импортированные группы",
    totalTypes: "Всего типов",
    importedTypes: "Импортированные типы",
    skippedAccountDetails: "Детали пропущенных аккаунтов",
    accountTitle: "Название аккаунта",
    accountName: "Имя аккаунта",
    accountId: "ID аккаунта",
    importFailed: "Импорт не удался",
    importFailedStatus: "Импорт не удался",
    unknownError: "Во время импорта произошла неизвестная ошибка",
    close: "Закрыть",
    cancel: "Отмена",
    previousStep: "Предыдущий шаг",
    startImport: "Начать импорт",
    nextStep: "Следующий шаг",
    refreshData: "Обновить данные",
    selectImportFileMessage: "Пожалуйста, выберите файл для импорта",
    enterBackupPasswordMessage: "Пожалуйста, введите пароль для распаковки",
    selectFileSuccess: "Файл для импорта успешно выбран",
    selectFileFailed: "Не удалось выбрать файл для импорта",
    preparingToImport: "Подготовка к импорту...",
    validatingFileAndPassword: "Проверка файла и пароля...",
    vaultImportSuccess: "Хранилище паролей успешно импортировано",
    dataRefreshed: "Данные обновлены",
    importInProgressWarning: "Импорт в процессе, пожалуйста, подождите...",
  },

  // Генератор паролей
  passwordGenerator: {
    title: "Генератор паролей",
    selectRule: "Выбрать правило пароля",
    selectRulePlaceholder: "Пожалуйста, выберите правило пароля",
    generalRule: "Общее правило пароля",
    customRule: "Пользовательское правило пароля",
    includeUppercase: "Включить заглавные буквы",
    includeLowercase: "Включить строчные буквы",
    includeNumbers: "Включить цифры",
    includeSpecialChars: "Включить специальные символы",
    passwordLength: "Длина пароля",
    customSpecialChars: "Пользовательские специальные символы",
    defaultSpecialChars:
      "Оставьте пустым для использования специальных символов по умолчанию",
    generatePassword: "Сгенерировать пароль",
    usePassword: "Использовать этот пароль",
    generatedPassword: "Сгенерированный пароль",
    clickToGenerate: "Нажмите кнопку генерации пароля",
    selectCharType: "Пожалуйста, выберите хотя бы один тип символов",
    enterPattern: "Пожалуйста, введите шаблон пароля",
    generateFirst: "Сначала сгенерируйте пароль",
  },

  // Новые переводы для настроек правил паролей
  passwordRuleSettings: {
    title: "Настроить правила генерации пароля",
    savedRules: "Сохраненные правила паролей",
    newRule: "Новое правило",
    ruleName: "Название правила",
    description: "Описание",
    operation: "Операция",
    edit: "Редактировать",
    delete: "Удалить",
    generalRule: "Общее правило пароля",
    includeUppercase: "Включить заглавные буквы",
    includeLowercase: "Включить строчные буквы",
    includeNumbers: "Включить цифры",
    includeSpecialChars: "Включить специальные символы",
    passwordLength: "Длина пароля",
    customSpecialChars: "Пользовательские специальные символы",
    customSpecialCharsPlaceholder:
      "Оставьте пустым для использования специальных символов по умолчанию: !@#$%^&*()_+-=[]{}|;:,.<>?",
    customRulesDescription: "Описание пользовательских правил паролей",
    close: "Закрыть",
    saveSettings: "Сохранить настройки",
    editRule: "Редактировать правило",
    newRuleDialog: "Новое правило",
    ruleNameRequired: "Пожалуйста, введите название правила",
    ruleDescriptionPlaceholder: "Пожалуйста, введите описание правила",
    ruleType: "Тип правила",
    general: "Общее правило",
    custom: "Пользовательское правило",
    passwordPattern: "Шаблон пароля",
    passwordPatternPlaceholder: "Пример: Aaa111",
    save: "Сохранить",
    cancel: "Отмена",
    confirmDeleteRule: 'Вы уверены, что хотите удалить правило "{ruleName}"?',
    deleteConfirmation: "Подтверждение удаления",
    confirmDelete: "Подтвердить",
    ruleDeletedSuccess: "Правило успешно удалено",
    ruleUpdateSuccess: "Правило успешно обновлено",
    ruleCreateSuccess: "Правило успешно создано",
    enterRuleName: "Пожалуйста, введите название правила",
    settingsSaved: "Настройки успешно сохранены",
  },

  // Новые переводы для результатов поиска
  searchResults: {
    groups: "Группы",
    accounts: "Аккаунты",
    noUrl: "Нет URL",
    belongsToGroup: "Принадлежит группе",
    unknownGroup: "Неизвестная группа",
    noSearchResults: "Результаты поиска не найдены",
    tryOtherKeywords: "Попробуйте другие ключевые слова",
  },

  // Новые переводы для строки состояния
  statusBar: {
    authInfo: "Информация об авторизации: Pro-версия активирована",
    usageCount: "Количество использований: {count} раз",
    usageDays: "Дней использования: {days} дней",
  },

  // Новые переводы для контекстного меню вкладок
  tabContextMenu: {
    rename: "Переименовать",
    deleteTab: "Удалить вкладку",
    newTab: "Новая вкладка",
    moveUp: "Переместить вверх",
    moveDown: "Переместить вниз",
    promptNewName: "Пожалуйста, введите новое имя вкладки",
    renameTab: "Переименовать вкладку",
    confirm: "Подтвердить",
    cancel: "Отмена",
    tabName: "Имя вкладки",
    tabNameCannotBeEmpty: "Имя вкладки не может быть пустым",
    tabNameNotChanged: "Имя вкладки не изменилось",
    userCanceledRename: "Пользователь отменил операцию переименования",
    confirmDeleteTab:
      'Вы уверены, что хотите удалить вкладку "{tabName}"?\nПосле удаления все аккаунты в этой вкладке также будут удалены, эта операция необратима.',
    deleteConfirmation: "Удалить вкладку",
    confirmDelete: "Подтвердить удаление",
    userCanceledDelete: "Пользователь отменил операцию удаления",
    promptNewTabName: "Пожалуйста, введите новое имя вкладки",
    userCanceledNewTab: "Пользователь отменил операцию создания новой вкладки",
  },

  // Новые переводы для боковой панели вкладок
  tabsSidebar: {
    emptyTabs: "Нет вкладок",
    createTabHint: "Нажмите кнопку ниже, чтобы создать вкладку",
    newTab: "Новая вкладка",
  },

  // Новые переводы для тестового диалога
  testDialog: {
    title: "Тестовый диалог",
    description1:
      "Это тестовый диалог для проверки правильности работы функции всплывающего окна.",
    description2:
      "Если вы видите этот диалог, это означает, что функция всплывающего окна работает правильно.",
    close: "Закрыть",
  },

  // Новые переводы для строки заголовка
  titleBar: {
    defaultAppTitle: "Менеджер паролей",
  },

  // Новые переводы для сервиса событий блокировки
  lockEventService: {
    logPrefix: "Служба событий блокировки",
    frontendStateUpdated: "Состояние интерфейса обновлено",
    redirectToLogin: "Перенаправление на страницу входа",
    vaultAutoLocked:
      "Хранилище паролей автоматически заблокировано, пожалуйста, войдите снова",
    sensitiveDataCleared: "Конфиденциальные данные очищены",
    clearSensitiveDataFailed: "Не удалось очистить конфиденциальные данные",
    windowEventListenersSet: "Установлены обработчики событий окна",
    windowLostFocus: "Окно потеряло фокус",
    windowGainedFocus: "Окно получило фокус",
    backendNotifiedMinimize: "Бэкенд уведомлен о сворачивании окна",
    notifyMinimizeFailed: "Не удалось уведомить о сворачивании окна",
    backendNotifiedFocus: "Бэкенд уведомлен о получении окном фокуса",
    notifyFocusFailed: "Не удалось уведомить о получении окном фокуса",
    manualLockCheckTriggered: "Запущена ручная проверка блокировки",
  },

  // Утилиты аккаунта
  accountUtils: {
    defaultAccountTitle: "Аккаунт",
    startCopyingPassword:
      "Начинаем копирование пароля, ID аккаунта: {accountId}",
    passwordCopiedSuccess:
      "Пароль скопирован в буфер обмена (автоочистка через 10 секунд)",
    passwordCopySuccess: "Пароль успешно скопирован, ID аккаунта: {accountId}",
    passwordCopyFailed: "Не удалось скопировать пароль",
    startGettingPassword: "Начинаем получение пароля, ID аккаунта: {accountId}",
    passwordGetSuccess: "Пароль успешно получен, ID аккаунта: {accountId}",
    passwordGetFailed: "Не удалось получить пароль",
    startGettingAccountDetail:
      "Начинаем получение деталей аккаунта, ID аккаунта: {accountId}",
    accountDetailGetSuccess:
      "Детали аккаунта успешно получены, ID аккаунта: {accountId}",
    accountDetailGetFailed: "Не удалось получить детали аккаунта",
  },
};
