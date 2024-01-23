export type Workbook = {
  id: number;
  name: string;
  description: string;
  createdAt: string;
  updatedAt: string;
};

export type Problem = {
  // id: number;
  type: string;
  properties: { [key: string]: string };
};

export type WorkbookWithProblems = {
  id: number;
  name: string;
  description: string;
  problems: Problem[];
};
