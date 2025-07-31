import { types } from '@models';

export interface ViewStateKey {
  viewName: string;
  tabName: types.DataFacet;
}

export const viewStateKeyToString = (key: ViewStateKey): string => {
  return `${key.viewName}/${key.tabName}/`;
};
