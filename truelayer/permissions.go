package truelayer

type Permission string

const PermissionAccounts Permission = "accounts"
const PermissionBalance Permission = "balance"
const PermissionCards Permission = "cards"
const PermissionTransactions Permission = "transactions"
const PermissionDirectDebits Permission = "direct_debits"
const PermissionStandingOrders Permission = "standing_orders"
const PermissionOfflineAccess Permission = "offline_access"
const PermissionInfo Permission = "info"
