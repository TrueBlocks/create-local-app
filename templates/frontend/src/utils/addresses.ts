export const getAddressString = (address: unknown): string => {
  if (typeof address === 'string') {
    return address;
  }

  if (address && typeof address === 'object' && 'address' in address) {
    const addrObj = address as { address?: string[] | string };
    if (Array.isArray(addrObj.address)) {
      return addrObj.address.join('');
    }
    if (typeof addrObj.address === 'string') {
      return addrObj.address;
    }
  }
  return '';
};

export const getDisplayAddress = (address: unknown): string => {
  const fullAddress = getAddressString(address);
  if (fullAddress.length > 10) {
    return `${fullAddress.slice(0, 6)}...${fullAddress.slice(-4)}`;
  }
  return fullAddress;
};
