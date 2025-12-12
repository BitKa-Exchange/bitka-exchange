import type { SidebarsConfig } from "@docusaurus/plugin-content-docs";

const sidebar: SidebarsConfig = {
  apisidebar: [
    {
      type: "doc",
      id: "openapi/bitka-exchange-api",
    },
    {
      type: "category",
      label: "Auth",
      items: [
        {
          type: "doc",
          id: "openapi/login",
          label: "Login",
          className: "api-method post",
        },
        {
          type: "doc",
          id: "openapi/register",
          label: "Register",
          className: "api-method post",
        },
        {
          type: "doc",
          id: "openapi/refresh-access-token",
          label: "Refresh access token",
          className: "api-method post",
        },
        {
          type: "doc",
          id: "openapi/logout",
          label: "Logout",
          className: "api-method post",
        },
        {
          type: "doc",
          id: "openapi/jwks-json-web-key-set",
          label: "JWKS (JSON Web Key Set)",
          className: "api-method get",
        },
      ],
    },
    {
      type: "category",
      label: "Users",
      items: [
        {
          type: "doc",
          id: "openapi/get-current-user-profile",
          label: "Get current user profile",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "openapi/update-current-user-profile",
          label: "Update current user profile",
          className: "api-method patch",
        },
        {
          type: "doc",
          id: "openapi/get-user-profile-by-id",
          label: "Get user profile by ID",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "openapi/update-current-user-profile",
          label: "Update current user profile",
          className: "api-method patch",
        },
        {
          type: "doc",
          id: "openapi/change-password-for-the-authenticated-user",
          label: "Change password for the authenticated user",
          className: "api-method post",
        },
      ],
    },
    {
      type: "category",
      label: "Ledger",
      items: [
        {
          type: "doc",
          id: "openapi/list-ledger-accounts",
          label: "List ledger accounts",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "openapi/get-ledger-account-by-id",
          label: "Get ledger account by id",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "openapi/list-transactions",
          label: "List transactions",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "openapi/create-a-transaction-debit-credit-transfer",
          label: "Create a transaction (debit/credit/transfer)",
          className: "api-method post",
        },
        {
          type: "doc",
          id: "openapi/get-transaction-by-id",
          label: "Get transaction by id",
          className: "api-method get",
        },
      ],
    },
    {
      type: "category",
      label: "Orders",
      items: [
        {
          type: "doc",
          id: "openapi/create-order",
          label: "Create order",
          className: "api-method post",
        },
        {
          type: "doc",
          id: "openapi/list-orders-filterable",
          label: "List orders (filterable)",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "openapi/get-order-by-id",
          label: "Get order by id",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "openapi/cancel-order-by-id",
          label: "Cancel order by id",
          className: "api-method delete",
        },
      ],
    },
    {
      type: "category",
      label: "MarketData",
      items: [
        {
          type: "doc",
          id: "openapi/list-available-symbols-market-pairs",
          label: "List available symbols (market pairs)",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "openapi/get-candles-ohlcv-for-a-symbol-and-interval-historical",
          label: "Get candles (OHLCV) for a symbol and interval (historical)",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "openapi/get-trades-ticks-for-a-symbol-historical",
          label: "Get trades (ticks) for a symbol (historical)",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "openapi/orderbook-snapshot-top-n-for-a-symbol",
          label: "Orderbook snapshot (top N) for a symbol",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "openapi/historical-orderbook-deltas-snapshots-for-reconciliation",
          label: "Historical orderbook deltas / snapshots for reconciliation",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "openapi/real-time-streaming-info-web-socket",
          label: "Real-time streaming info (WebSocket)",
          className: "api-method get",
        },
      ],
    },
  ],
};

export default sidebar.apisidebar;
