/**
 * 繁體中文語言包
 * @author 陳鳳慶
 * @date 2025-10-05
 */

export default {
  // 通用
  common: {
    confirm: "確認",
    cancel: "取消",
    save: "儲存",
    delete: "刪除",
    edit: "編輯",
    add: "新增",
    search: "搜尋",
    close: "關閉",
    ok: "確定",
    yes: "是",
    no: "否",
    loading: "載入中...",
    success: "成功",
    error: "錯誤",
    warning: "警告",
    info: "資訊",
    copy: "複製",
    paste: "貼上",
    cut: "剪下",
    refresh: "重新整理",
    reset: "重設",
    clear: "清空",
    back: "返回",
    next: "下一步",
    previous: "上一步",
    finish: "完成",
    submit: "提交",
    import: "匯入",
    export: "匯出",
    settings: "設定",
    help: "說明",
    about: "關於",
    more: "更多",
  },

  // 登入頁面
  login: {
    title: "密碼庫登入",
    password: "密碼",
    passwordPlaceholder: "請輸入密碼庫密碼",
    login: "登入",
    createVault: "建立密碼庫",
    openVault: "開啟其他密碼庫",
    invalidPassword: "密碼錯誤",
    loginSuccess: "登入成功",
    loginFailed: "登入失敗",
    recentVaults: "最近使用",
    passwordHint1: "登入密碼是密碼庫加密和解密的關鍵金鑰",
    passwordHint2: "請妥善保管您的密碼，一旦遺失無法找回",
    passwordHint3: "登入密碼以加密形式儲存在資料庫中",
    vaultFile: "密碼庫檔案",
    vaultFilePlaceholder: "請選擇密碼庫檔案",
    vaultPathRequired: "請輸入密碼庫檔案路徑",
    passwordRequired: "請輸入登入密碼",
    passwordMinLength: "密碼長度不能少於6位",
    selectFileFailed: "選擇檔案失敗",
    passwordEmpty: "密碼不能為空",
    passwordMinLength8: "密碼長度至少8位",
    passwordNeedUppercase: "密碼必須包含大寫字母",
    passwordNeedLowercase: "密碼必須包含小寫字母",
    passwordNeedNumber: "密碼必須包含數字",
    passwordNeedSpecialChar: "密碼必須包含特殊字元(!@#$%^&*等)",
    passwordStrengthOk: "密碼強度符合要求",
    vaultFileNotExists: "密碼庫檔案不存在",
    createVaultPrompt:
      "請輸入密碼庫名稱（無需副檔名）\n\n系統將自動在目前程式所在目錄下建立data目錄，並將密碼庫檔案儲存到data目錄下，根據名稱自動建立.db檔案。",
    vaultNamePlaceholder: "例如: my_vault",
    setPasswordPrompt:
      "請設定登入密碼\n\n密碼要求：\n• 至少8位字元\n• 包含大寫字母\n• 包含小寫字母\n• 包含數字\n• 包含特殊字元(!@#$%^&*等)",
    setPassword: "設定密碼",
    passwordRequirements: "密碼要求",
    passwordReq1: "至少8位字元",
    passwordReq2: "包含大寫字母",
    passwordReq3: "包含小寫字母",
    passwordReq4: "包含數字",
    passwordReq5: "包含特殊字元(!@#$%^&*等)",
    createVaultSuccess: "密碼庫建立成功",
    createVaultFailed: "建立密碼庫失敗",
    switchToFullMode: "已切換到完整模式，請選擇密碼庫檔案",
    switchToCreateMode: "已切換到完整模式，請建立新的密碼庫",
  },

  // 主介面
  main: {
    title: "wepass - 密碼管理器",
    searchPlaceholder: "搜尋帳號...",
    noData: "暫無資料",
    noSearchResults: "未找到符合的結果",
    searchResults: "搜尋結果",
    groups: "分組",
    group: "分組",
    accounts: "帳號",
    allAccounts: "所有帳號",
    noAddress: "無地址",
    belongsToGroup: "所屬分組",
    inputUsernameAndPassword: "輸入使用者名稱和密碼",
    inputUsername: "輸入使用者名稱",
    inputPassword: "輸入密碼",
    tryOtherKeywords: "請嘗試其他關鍵字",
    unknownGroup: "未知分組",
  },

  // 設定對話框
  settings: {
    title: "設定",
    general: "一般",
    log: "日誌",
    lock: "鎖定",
    theme: "主題",
    language: "語言",
    lightTheme: "淺色",
    darkTheme: "深色",
    selectTheme: "選擇主題",
    selectLanguage: "選擇語言",
    logSettings: "日誌設定",
    enableInfoLog: "啟用資訊日誌",
    enableDebugLog: "啟用除錯日誌",
    lockSettings: "鎖定設定",
    enableAutoLock: "啟用自動鎖定",
    enableTimerLock: "啟用定時鎖定",
    enableMinimizeLock: "啟用最小化鎖定",
    lockTimeMinutes: "鎖定時間（分鐘）",
    enableSystemLock: "啟用系統鎖定",
    systemLockMinutes: "系統鎖定時間（分鐘）",
    settingsSaved: "設定已儲存",
    loadSettingsFailed: "載入設定失敗",
    saveSettingsFailed: "儲存設定失敗",
  },

  // 帳號管理
  account: {
    title: "標題",
    username: "使用者名稱",
    password: "密碼",
    url: "網址",
    type: "類型",
    notes: "備註",
    icon: "圖示",
    favorite: "收藏",
    useCount: "使用次數",
    lastUsed: "最後使用",
    created: "建立時間",
    updated: "更新時間",
    inputMethod: "輸入方式",
    addAccount: "新增帳號",
    editAccount: "編輯帳號",
    deleteAccount: "刪除帳號",
    copyUsername: "複製使用者名稱",
    copyPassword: "複製密碼",
    openUrl: "開啟網址",
    generatePassword: "產生密碼",
    showPassword: "顯示密碼",
    hidePassword: "隱藏密碼",
    accountSaved: "帳號已儲存",
    accountDeleted: "帳號已刪除",
    confirmDelete: "確認刪除此帳號嗎？",
    deleteConfirmTitle: "刪除確認",
    copiedToClipboard: "已複製到剪貼簿",
    copyFailed: "複製失敗",
    noAccounts: "暫無帳號",
    noAccountsHint: "點擊右上方「建立帳號」按鈕新增",
  },

  // 分組和類型
  group: {
    allGroups: "所有分組",
    addGroup: "新增分組",
    editGroup: "編輯分組",
    deleteGroup: "刪除分組",
    groupName: "分組名稱",
    groupIcon: "分組圖示",
    rename: "重新命名", // TODO: Translate
    createGroup: "建立群組", // TODO: Translate
    moveLeft: "向左移動", // TODO: Translate
    moveRight: "向右移動", // TODO: Translate
    newGroupName: "新群組名稱", // TODO: Translate
    renameGroup: "重新命名群組", // TODO: Translate
    newGroupPlaceholder: "群組名稱", // TODO: Translate
    groupNameCannotBeEmpty: "群組名稱不能為空", // TODO: Translate
    groupNameNotChanged: "群組名稱未變更", // TODO: Translate
    defaultGroupCannotBeDeleted: "預設群組無法刪除", // TODO: Translate
    confirmDeleteGroup:
      '確定要刪除群組 "{groupName}" 嗎？\n此群組下的所有標籤和帳號也將被刪除，此操作無法復原。', // TODO: Translate
    deleteGroupTitle: "刪除群組", // TODO: Translate
    confirmDelete: "確認刪除", // TODO: Translate
    defaultGroupName: "預設", // TODO: Translate
  },

  type: {
    allTypes: "所有類型",
    addType: "新增類型",
    editType: "編輯類型",
    deleteType: "刪除類型",
    typeName: "類型名稱",
    typeIcon: "類型圖示",
  },

  // 說明對話框
  help: {
    title: "說明",
    quickStart: "快速入門",
    welcome: "歡迎使用 WePassword",
    welcomeTitle: "歡迎使用 WePassword",
    description:
      "WePassword 是一個安全的密碼管理工具，幫助您安全地儲存和管理密碼。",
    welcomeDescription:
      "WePassword 是一個安全的密碼管理工具，幫助您安全地儲存和管理密碼。",
    basicOperations: "基本操作",
    createVault: "建立密碼庫：首次使用時，系統會引導您建立一個新的密碼庫",
    addPassword: "新增密碼：點擊「新增」按鈕，輸入網站、使用者名稱和密碼資訊",
    searchPassword: "搜尋密碼：使用搜尋框快速查找特定的密碼條目",
    editPassword: "編輯密碼：雙擊密碼條目或點擊編輯按鈕進行修改",
    deletePassword: "刪除密碼：選擇密碼條目後點擊刪除按鈕",
    features: "功能說明", // TODO: Translate
    mainFeatures: "主要功能", // TODO: Translate
    passwordGenerator: "密碼產生器", // TODO: Translate
    passwordGeneratorDesc: "產生強密碼，支援自訂規則", // TODO: Translate
    groupManagement: "分組管理", // TODO: Translate
    groupManagementDesc: "使用標籤頁對密碼進行分類管理", // TODO: Translate
    secureEncryption: "安全加密", // TODO: Translate
    secureEncryptionDesc: "所有密碼資料都經過加密儲存", // TODO: Translate
    importExport: "匯入/匯出", // TODO: Translate
    importExportDesc: "支援從其他密碼管理器匯入資料", // TODO: Translate
    backupRestore: "備份/還原", // TODO: Translate
    backupRestoreDesc: "定期備份您的密碼庫", // TODO: Translate
    passwordGenerationRules: "密碼產生規則", // TODO: Translate
    passwordGenerationRulesDesc: "支援多種字元集和自訂規則：", // TODO: Translate
    lowercaseLetters: "小寫字母", // TODO: Translate
    mixedCaseLetters: "大小寫混合字母", // TODO: Translate
    uppercaseLetters: "大寫字母", // TODO: Translate
    digits: "數字", // TODO: Translate
    specialCharacters: "特殊字元", // TODO: Translate
    customCharacterSet: "自訂字元集", // TODO: Translate
    securityTips: "安全提示", // TODO: Translate
    securitySuggestions: "安全建議", // TODO: Translate
    masterPassword: "主密碼", // TODO: Translate
    masterPasswordDesc: "設定一個強主密碼，並牢記它", // TODO: Translate
    regularBackup: "定期備份", // TODO: Translate
    regularBackupDesc: "定期備份您的密碼庫檔案", // TODO: Translate
    timelyUpdate: "及時更新", // TODO: Translate
    timelyUpdateDesc: "定期更新重要帳戶的密碼", // TODO: Translate
    avoidRepetition: "避免重複", // TODO: Translate
    avoidRepetitionDesc: "不要在多個網站使用相同密碼", // TODO: Translate
    safeEnvironment: "安全環境", // TODO: Translate
    safeEnvironmentDesc: "在安全的環境中使用密碼管理器", // TODO: Translate
    precautions: "注意事項", // TODO: Translate
    precaution1: "請妥善保管您的主密碼，忘記主密碼將無法還原資料。", // TODO: Translate
    precaution2: "建議定期變更主密碼以提高安全性。", // TODO: Translate
    precaution3: "不要在公用電腦上儲存密碼庫檔案。", // TODO: Translate
    faq: "常見問題", // TODO: Translate
    faqTitle: "常見問題解答", // TODO: Translate
    faq1_q: "忘記主密碼怎麼辦？", // TODO: Translate
    faq1_a:
      "很遺憾，如果忘記主密碼，將無法還原密碼庫中的資料。建議您定期備份密碼庫，並將主密碼記錄在安全的地方。", // TODO: Translate
    faq2_q: "如何備份密碼庫？", // TODO: Translate
    faq2_a:
      "您可以複製密碼庫檔案到安全的位置，或使用「更多」選單中的匯出功能。", // TODO: Translate
    faq3_q: "支援哪些檔案格式？", // TODO: Translate
    faq3_a: "目前支援 .db 和 .vault 格式的密碼庫檔案。", // TODO: Translate
    faq4_q: "如何產生安全的密碼？", // TODO: Translate
    faq4_a:
      "使用「更多」選單中的「產生密碼」功能，建議使用包含大小寫字母、數字和特殊字元的組合，長度至少12位。", // TODO: Translate
  },

  // 更多選單
  moreMenu: {
    more: "更多",
    selectNewVault: "選擇新的密碼庫",
    openVaultDirectory: "開啟密碼庫目錄",
    generatePassword: "產生密碼",
    setPasswordRules: "設定產生密碼規則",
    changeLoginPassword: "修改登入密碼",
    exportVault: "匯出密碼庫",
    importVault: "匯入密碼庫",
    changeLog: "更新日誌",
    settings: "設定",
    help: "說明",
    about: "關於",
    lockVault: "鎖定密碼庫",
    logout: "登出",
    oldLoginPasswordLabel: "舊登入密碼",
    oldLoginPasswordPlaceholder: "請輸入目前的登入密碼",
    newLoginPasswordLabel: "新登入密碼",
    newLoginPasswordPlaceholder: "請輸入新的登入密碼",
    confirmNewPasswordLabel: "確認新登入密碼",
    confirmNewPasswordPlaceholder: "請再次輸入新的登入密碼",
    selectVaultFile: "選擇密碼庫檔案",
    selectVaultFilePrompt: "請選擇要開啟的密碼庫檔案：",
    selectVaultFilePlaceholder: "請選擇密碼庫檔案路徑",
    browse: "瀏覽",
    supportedFileFormats: "支援的檔案格式：.db, .vault",
    openFileFailed: "開啟檔案失敗",
    passwordCopied: "密碼已複製到剪貼簿",
    copyFailedManual: "複製失敗，請手動複製密碼",
    passwordChangeInProgress: "密碼變更進行中，請稍候...",
    formRefNotFound: "找不到表單引用",
    oldPasswordIncorrect: "舊密碼不正確，請重新輸入",
    changePasswordConfirm:
      "變更登入密碼將重新加密所有帳號資料。此操作不可逆，請確保您已備份目前密碼庫。是否繼續？",
    passwordChangeSuccess: "登入密碼變更成功！",
    passwordChangeFailed: "登入密碼變更失敗",
    importSuccessDataUpdated: "匯入成功，資料已更新",
    importSuccessDataRefreshFailed:
      "匯入成功，但資料更新失敗，請手動重新整理頁面",
    vaultDirectoryOpened: "密碼庫目錄已開啟",
    openDirectoryFailed: "開啟目錄失敗",
    lockVaultConfirm: "您確定要鎖定密碼庫嗎？鎖定後需要重新輸入密碼才能存取。",
    vaultLocked: "密碼庫已鎖定",
    lockVaultFailed: "鎖定密碼庫失敗",
    logoutConfirm: "您確定要登出嗎？登出後需要重新輸入密碼才能存取密碼庫。",
    logoutSuccess: "登出成功",
    oldLoginPasswordRequired: "請輸入舊密碼",
    newLoginPasswordRequired: "請輸入新密碼",
    newLoginPasswordMinLength: "密碼長度至少8位",
    newLoginPasswordStrength: "密碼必須包含大寫字母、小寫字母、數字和特殊字元",
    confirmNewLoginPasswordRequired: "請確認新密碼",
    passwordsDoNotMatch: "密碼不符",
    selectNewVaultConfirm:
      "選擇新的密碼庫會將您帶回登入畫面，您需要重新選擇密碼庫檔案並輸入密碼。是否繼續？",
    selectNewVaultSuccess: "您已返回登入畫面，請選擇新的密碼庫",
    pleaseSelectFile: "請先選擇檔案",
    openingFile: "正在開啟檔案：",
    open: "開啟",
    continue: "繼續",
    verifying: "驗證中...",
    changingPassword: "變更中...",
  },

  // 右鍵選單
  contextMenu: {
    inputUsernameAndPassword: "輸入使用者名稱和密碼", // TODO: Translate
    openUrl: "開啟地址",
    duplicate: "產生副本",
    view: "檢視",
    edit: "修改",
    changeGroup: "更改分組",
    copyUsername: "複製帳號",
    copyPassword: "複製密碼",
    copyUsernameAndPassword: "複製帳號密碼",
    copyUrl: "複製地址",
    copyTitle: "複製標題",
    copyNotes: "複製備註",
    delete: "刪除",
  },

  // 匯出功能
  export: {
    title: "匯出密碼庫",
    steps: {
      verifyPassword: "驗證密碼",
      selectAccounts: "選擇帳號",
      setBackup: "設定備份",
      exportComplete: "匯出完成",
    },
    verifyPasswordTitle: "驗證登入密碼",
    verifyPasswordDesc: "請輸入目前密碼庫的登入密碼以繼續匯出操作。",
    loginPassword: "登入密碼",
    loginPasswordPlaceholder: "請輸入登入密碼",
    selectAccountsTitle: "選擇要匯出的帳號",
    exportAll: "全部匯出",
    exportByGroup: "按分組匯出",
    exportByType: "按類型匯出", // TODO: Translate
    exportSelected: "選擇匯出",
    selectGroups: "選擇分組", // TODO: Translate
    selectAll: "全選", // TODO: Translate
    clearAll: "取消全選", // TODO: Translate
    groupSelectionSummary:
      "已選 {count} 個群組，預計匯出 {accountCount} 個帳號", // TODO: Translate
    selectTypes: "選擇類型", // TODO: Translate
    typeSelectionSummary: "已選 {count} 個類型，預計匯出 {accountCount} 個帳號", // TODO: Translate
    selectAccounts: "選擇帳號", // TODO: Translate
    loadingAccounts: "正在載入帳號...", // TODO: Translate
    noAccounts: "沒有帳號資料", // TODO: Translate
    setBackupTitle: "設定備份密碼和匯出路徑", // TODO: Translate
    backupPassword: "備份密碼", // TODO: Translate
    backupPasswordPlaceholder: "請輸入備份密碼", // TODO: Translate
    generate: "產生", // TODO: Translate
    backupPasswordTip: "備份密碼用於加密匯出的資料，請妥善保管", // TODO: Translate
    exportPath: "匯出路徑", // TODO: Translate
    exportPathPlaceholder: "請選擇匯出路徑", // TODO: Translate
    browse: "瀏覽", // TODO: Translate
    exporting: "正在匯出密碼庫...", // TODO: Translate
    exportSuccessTitle: "匯出成功", // TODO: Translate
    exportSuccessSubTitle: "密碼庫已成功匯出至：{path}", // TODO: Translate
    exportFailedTitle: "匯出失敗", // TODO: Translate
    startExport: "開始匯出", // TODO: Translate
    openFolder: "開啟資料夾", // TODO: Translate
    loginPasswordRequired: "請輸入登入密碼", // TODO: Translate
    backupPasswordRequired: "請輸入備份密碼", // TODO: Translate
    backupPasswordMinLength: "備份密碼長度不能少於6位", // TODO: Translate
    exportPathRequired: "請選擇匯出路徑", // TODO: Translate
  },

  // 匯入功能
  import: {
    title: "匯入密碼庫",
    selectFile: "選擇檔案",
    selectFileDesc: "請選擇要匯入的密碼庫檔案",
    fileFormat: "檔案格式",
    importProgress: "匯入進度",
    importComplete: "匯入完成",
  },

  // 新增的匯入功能翻譯
  importVault: {
    title: "匯入密碼庫",
    step1: "選擇檔案",
    step2: "驗證密碼",
    step3: "匯入完成",
    step1Title: "選擇匯入檔案",
    step1Description: "請選擇要匯入的密碼庫備份檔案（ZIP格式）。",
    importFile: "匯入檔案",
    selectImportFilePlaceholder: "請選擇匯入檔案",
    browse: "瀏覽",
    step2Title: "輸入解壓縮密碼",
    step2Description: "請輸入備份檔案的解壓縮密碼（備份密碼）。",
    backupPassword: "解壓縮密碼",
    enterBackupPasswordPlaceholder: "請輸入解壓縮密碼",
    importingVault: "正在匯入密碼庫...",
    importComplete: "匯入完成",
    importSuccess: "匯入成功",
    vaultImportedSuccessfully: "密碼庫已成功匯入",
    importReport: "匯入報告",
    totalAccounts: "總帳號數",
    successfullyImported: "成功匯入",
    skippedAccounts: "跳過帳號",
    errorAccounts: "錯誤帳號",
    totalGroups: "總分組數",
    importedGroups: "匯入分組",
    totalTypes: "總類型數",
    importedTypes: "匯入類型",
    skippedAccountDetails: "跳過的帳號詳情",
    accountTitle: "帳號標題",
    accountName: "帳號名稱",
    accountId: "帳號ID",
    importFailed: "匯入失敗",
    importFailedStatus: "匯入失敗",
    unknownError: "匯入過程中發生未知錯誤",
    close: "關閉",
    cancel: "取消",
    previousStep: "上一步",
    startImport: "開始匯入",
    nextStep: "下一步",
    refreshData: "重新整理資料",
    selectImportFileMessage: "請選擇匯入檔案",
    enterBackupPasswordMessage: "請輸入解壓縮密碼",
    selectFileSuccess: "匯入檔案選擇成功",
    selectFileFailed: "選擇匯入檔案失敗",
    preparingToImport: "準備匯入...",
    validatingFileAndPassword: "驗證檔案和密碼...",
    vaultImportSuccess: "密碼庫匯入成功",
    dataRefreshed: "資料已重新整理",
    importInProgressWarning: "匯入正在進行中，請稍候...",
  },

  // 密碼產生器
  passwordGenerator: {
    title: "產生密碼",
    selectRule: "選擇密碼規則",
    selectRulePlaceholder: "請選擇密碼規則",
    generalRule: "通用密碼規則",
    customRule: "自訂密碼規則",
    includeUppercase: "包含大寫字母",
    includeLowercase: "包含小寫字母",
    includeNumbers: "包含數字",
    includeSpecialChars: "包含特殊字元",
    passwordLength: "密碼長度",
    customSpecialChars: "自訂特殊字元",
    defaultSpecialChars: "留空使用預設特殊字元",
    generatePassword: "產生密碼",
    usePassword: "使用此密碼",
    generatedPassword: "產生的密碼",
    clickToGenerate: "點擊產生密碼按鈕",
    selectCharType: "請至少選擇一種字元類型",
    enterPattern: "請輸入密碼模式",
    generateFirst: "請先產生密碼",
  },

  // 新增的密碼規則設定翻譯
  passwordRuleSettings: {
    title: "設定產生密碼規則",
    savedRules: "已儲存的密碼規則",
    newRule: "新增規則",
    ruleName: "規則名稱",
    description: "描述",
    operation: "操作",
    edit: "編輯",
    delete: "刪除",
    generalRule: "通用密碼規則",
    includeUppercase: "包含大寫字母",
    includeLowercase: "包含小寫字母",
    includeNumbers: "包含數字",
    includeSpecialChars: "包含特殊字元",
    passwordLength: "密碼長度",
    customSpecialChars: "自訂特殊字元",
    customSpecialCharsPlaceholder:
      "留空使用預設特殊字元: !@#$%^&*()_+-=[]{}|;:,.<>?",
    customRulesDescription: "自訂密碼規則說明",
    close: "關閉",
    saveSettings: "儲存設定",
    editRule: "編輯規則",
    newRuleDialog: "新增規則",
    ruleNameRequired: "請輸入規則名稱",
    ruleDescriptionPlaceholder: "請輸入規則描述",
    ruleType: "規則類型",
    general: "通用規則",
    custom: "自訂規則",
    passwordPattern: "密碼模式",
    passwordPatternPlaceholder: "例如：Aaa111",
    save: "儲存",
    cancel: "取消",
    confirmDeleteRule: "確定要刪除規則「{ruleName}」嗎？",
    deleteConfirmation: "刪除確認",
    confirmDelete: "確定",
    ruleDeletedSuccess: "規則刪除成功",
    ruleUpdateSuccess: "規則更新成功",
    ruleCreateSuccess: "規則建立成功",
    enterRuleName: "請輸入規則名稱",
    settingsSaved: "設定儲存成功",
  },

  // 新增的搜尋結果翻譯
  searchResults: {
    groups: "分組",
    accounts: "帳號",
    noUrl: "無網址",
    belongsToGroup: "所屬分組",
    unknownGroup: "未知分組",
    noSearchResults: "未找到搜尋結果",
    tryOtherKeywords: "請嘗試其他關鍵字",
  },

  // 新增的狀態列翻譯
  statusBar: {
    authInfo: "授權資訊: 已啟用專業版",
    usageCount: "使用次數: {count}次",
    usageDays: "使用天數: {days}天",
  },

  // 新增的標籤上下文選單翻譯
  tabContextMenu: {
    rename: "重新命名",
    deleteTab: "刪除標籤",
    newTab: "新增標籤",
    moveUp: "上移",
    moveDown: "下移",
    promptNewName: "請輸入新的標籤名稱",
    renameTab: "重新命名標籤",
    confirm: "確定",
    cancel: "取消",
    tabName: "標籤名稱",
    tabNameCannotBeEmpty: "標籤名稱不能為空",
    tabNameNotChanged: "標籤名稱未變更",
    userCanceledRename: "使用者取消重新命名操作",
    confirmDeleteTab:
      "確定要刪除標籤「{tabName}」嗎？\n刪除後該標籤下的所有帳號也將被刪除，此操作不可恢復。",
    deleteConfirmation: "刪除標籤",
    confirmDelete: "確定刪除",
    userCanceledDelete: "使用者取消刪除操作",
    promptNewTabName: "請輸入新的標籤名稱",
    userCanceledNewTab: "使用者取消新增標籤操作",
  },

  // 新增的標籤側邊欄翻譯
  tabsSidebar: {
    emptyTabs: "暫無標籤",
    createTabHint: "點擊下方按鈕建立標籤",
    newTab: "新增標籤",
  },

  // 新增的測試對話框翻譯
  testDialog: {
    title: "測試對話框",
    description1: "這是一個測試對話框，用於驗證彈出視窗功能是否正常工作。",
    description2: "如果您能看到這個對話框，說明彈出視窗功能正常。",
    close: "關閉",
  },

  // 新增的標題列翻譯
  titleBar: {
    defaultAppTitle: "密碼管理器",
  },

  // 新增的鎖定事件服務翻譯
  lockEventService: {
    logPrefix: "鎖定事件服務",
    frontendStateUpdated: "前端狀態已更新",
    redirectToLogin: "跳轉到登入頁面",
    vaultAutoLocked: "密碼庫已自動鎖定，請重新登入",
    sensitiveDataCleared: "敏感資料已清理",
    clearSensitiveDataFailed: "清理敏感資料失敗",
    windowEventListenersSet: "視窗事件監聽器已設定",
    windowLostFocus: "視窗失去焦點",
    windowGainedFocus: "視窗獲得焦點",
    backendNotifiedMinimize: "已通知後端視窗最小化",
    notifyMinimizeFailed: "通知視窗最小化失敗",
    backendNotifiedFocus: "已通知後端視窗獲得焦點",
    notifyFocusFailed: "通知視窗獲得焦點失敗",
    manualLockCheckTriggered: "手動觸發鎖定檢查",
  },

  // 新增的帳號工具函數翻譯
  accountUtils: {
    defaultAccountTitle: "帳號",
    startCopyingPassword: "開始複製密碼，帳號ID: {accountId}",
    passwordCopiedSuccess: "密碼已複製到剪貼簿（10秒後自動清理）",
    passwordCopySuccess: "密碼複製成功，帳號ID: {accountId}",
    passwordCopyFailed: "複製密碼失敗",
    startGettingPassword: "開始獲取密碼，帳號ID: {accountId}",
    passwordGetSuccess: "密碼獲取成功，帳號ID: {accountId}",
    passwordGetFailed: "獲取密碼失敗",
    startGettingAccountDetail: "開始獲取帳號詳情，帳號ID: {accountId}",
    accountDetailGetSuccess: "帳號詳情獲取成功，帳號ID: {accountId}",
    accountDetailGetFailed: "獲取帳號詳情失敗",
  },
};
