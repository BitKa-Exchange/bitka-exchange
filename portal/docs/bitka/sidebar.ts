import type { SidebarsConfig } from "@docusaurus/plugin-content-docs";

const sidebar: SidebarsConfig = {
  apisidebar: [
    {
      type: "doc",
      id: "bitka/bitka-exchange-api",
    },
    {
      type: "category",
      label: "Auth",
      items: [
        {
          type: "doc",
          id: "bitka/login",
          label: "Login",
          className: "api-method post",
        },
        {
          type: "doc",
          id: "bitka/register",
          label: "Register",
          className: "api-method post",
        },
        {
          type: "doc",
          id: "bitka/refresh-access-token",
          label: "Refresh access token",
          className: "api-method post",
        },
        {
          type: "doc",
          id: "bitka/logout",
          label: "Logout",
          className: "api-method post",
        },
        {
          type: "doc",
          id: "bitka/jwks-json-web-key-set",
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
          id: "bitka/get-current-user-profile",
          label: "Get current user profile",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "bitka/update-current-user-profile",
          label: "Update current user profile",
          className: "api-method patch",
        },
        {
          type: "doc",
          id: "bitka/get-user-profile-by-id",
          label: "Get user profile by ID",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "bitka/update-current-user-profile",
          label: "Update current user profile",
          className: "api-method patch",
        },
        {
          type: "doc",
          id: "bitka/change-password-for-the-authenticated-user",
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
          id: "bitka/list-ledger-accounts",
          label: "List ledger accounts",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "bitka/get-ledger-account-by-id",
          label: "Get ledger account by id",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "bitka/list-transactions",
          label: "List transactions",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "bitka/create-a-transaction-debit-credit-transfer",
          label: "Create a transaction (debit/credit/transfer)",
          className: "api-method post",
        },
        {
          type: "doc",
          id: "bitka/get-transaction-by-id",
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
          id: "bitka/create-order",
          label: "Create order",
          className: "api-method post",
        },
        {
          type: "doc",
          id: "bitka/list-orders-filterable",
          label: "List orders (filterable)",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "bitka/get-order-by-id",
          label: "Get order by id",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "bitka/cancel-order-by-id",
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
          id: "bitka/list-available-symbols-market-pairs",
          label: "List available symbols (market pairs)",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "bitka/get-candles-ohlcv-for-a-symbol-and-interval-historical",
          label: "Get candles (OHLCV) for a symbol and interval (historical)",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "bitka/get-trades-ticks-for-a-symbol-historical",
          label: "Get trades (ticks) for a symbol (historical)",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "bitka/orderbook-snapshot-top-n-for-a-symbol",
          label: "Orderbook snapshot (top N) for a symbol",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "bitka/historical-orderbook-deltas-snapshots-for-reconciliation",
          label: "Historical orderbook deltas / snapshots for reconciliation",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "bitka/real-time-streaming-info-web-socket",
          label: "Real-time streaming info (WebSocket)",
          className: "api-method get",
        },
      ],
    },
  ],
};

export default sidebar.apisidebar;
