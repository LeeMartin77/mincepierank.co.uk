export type Maker = {
  id: string;
  name: string;
};

export type MakerPie = {
  makerid: string;
  id: string;
  displayname: string;
  fresh: boolean;
  labels: string[];
};

export type MakerPieRanking = {
  makerid: string;
  pieid: string;
  userid: string;
  pastry: number;
  filling: number;
  topping: number;
  looks: number;
  value: number;
  notes: string;
};

export enum StorageError {
  GenericError,
  NotFound,
}
