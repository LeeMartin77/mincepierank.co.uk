export type Maker = {
  id: string;
  name: string;
};

export type MakerPie = {
  makerId: string;
  id: string;
  displayName: string;
  fresh: boolean;
  labels: string[];
};

export type MakerPieRanking = {
  pieId: string;
  userId: string;
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
