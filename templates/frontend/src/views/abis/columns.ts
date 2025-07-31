// Copyright 2016, 2026 The {{ORG_NAME}} Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */
import { FormField } from '@components';
import { types } from '@models';

// EXISTING_CODE
// EXISTING_CODE

// Column configurations for the Abis data facets

export const getColumns = (dataFacet: types.DataFacet): FormField[] => {
  switch (dataFacet) {
    case types.DataFacet.DOWNLOADED:
      return getAbisColumns();
    case types.DataFacet.KNOWN:
      return getAbisColumns();
    case types.DataFacet.FUNCTIONS:
      return getFunctionsColumns();
    case types.DataFacet.EVENTS:
      return getFunctionsColumns();
    default:
      return [];
  }
};

const getAbisColumns = (): FormField[] => [
  // EXISTING_CODE
  // EXISTING_CODE
  {
    key: 'address',
    name: 'address',
    header: 'Address',
    label: 'Address',
    type: 'address',
    width: '340px',
    readOnly: true,
  },
  {
    key: 'name',
    name: 'name',
    header: 'Name',
    label: 'Name',
    type: 'text',
    width: '200px',
  },
  {
    key: 'fileSize',
    name: 'fileSize',
    header: 'File Size',
    label: 'File Size',
    type: 'number',
    width: '120px',
  },
  {
    key: 'nFunctions',
    name: 'nFunctions',
    header: 'Functions',
    label: 'Functions',
    type: 'number',
    width: '120px',
  },
  {
    key: 'nEvents',
    name: 'nEvents',
    header: 'Events',
    label: 'Events',
    type: 'number',
    width: '120px',
  },
  {
    key: 'actions',
    name: 'actions',
    header: 'Actions',
    label: 'Actions',
    editable: false,
    visible: true,
    type: 'button',
    width: '80px',
  },
];

const getFunctionsColumns = (): FormField[] => [
  // EXISTING_CODE
  // EXISTING_CODE
  {
    key: 'encoding',
    name: 'encoding',
    header: 'Encoding',
    label: 'Encoding',
    type: 'text',
    width: '200px',
  },
  {
    key: 'name',
    name: 'name',
    header: 'Name',
    label: 'Name',
    type: 'text',
    width: '200px',
  },
  {
    key: 'type',
    name: 'type',
    header: 'Type',
    label: 'Type',
    type: 'text',
    width: '200px',
  },
  {
    key: 'signature',
    name: 'signature',
    header: 'Signature',
    label: 'Signature',
    type: 'text',
    width: '200px',
  },
  {
    key: 'actions',
    name: 'actions',
    header: 'Actions',
    label: 'Actions',
    editable: false,
    visible: true,
    type: 'button',
    width: '80px',
  },
];

// EXISTING_CODE
// EXISTING_CODE
