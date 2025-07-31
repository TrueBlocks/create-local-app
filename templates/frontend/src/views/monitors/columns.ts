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

// Column configurations for the Monitors data facets

export const getColumns = (dataFacet: types.DataFacet): FormField[] => {
  switch (dataFacet) {
    case types.DataFacet.MONITORS:
      return getMonitorsColumns();
    default:
      return [];
  }
};

const getMonitorsColumns = (): FormField[] => [
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
    key: 'nRecords',
    name: 'nRecords',
    header: 'Records',
    label: 'Records',
    type: 'number',
    width: '120px',
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
    key: 'isEmpty',
    name: 'isEmpty',
    header: 'Empty',
    label: 'Empty',
    type: 'checkbox',
    width: '80px',
  },
  {
    key: 'lastScanned',
    name: 'lastScanned',
    header: 'Last Scanned',
    label: 'Last Scanned',
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

// EXISTING_CODE
// EXISTING_CODE
