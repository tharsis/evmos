local config = import 'default.jsonnet';

config {
  'eidon-chain_9002-1'+: {
    genesis+: {
      app_state+: {
        feemarket+: {
          params+: {
            min_gas_price: '0',
            no_base_fee: true,
            base_fee: '0',            
          },
        },
      },
    },
  },
}
