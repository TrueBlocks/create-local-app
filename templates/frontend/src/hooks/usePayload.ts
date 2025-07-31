import { useCallback } from 'react';

import { types } from '@models';

import { useActiveProject } from './useActiveProject';

export const usePayload = () => {
  const { lastAddress, lastChain } = useActiveProject();
  return useCallback(
    (dataFacet: types.DataFacet, address?: string) => {
      return types.Payload.createFrom({
        dataFacet,
        chain: lastChain,
        address: address || lastAddress,
      });
    },
    [lastChain, lastAddress],
  );
};
