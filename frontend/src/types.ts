export type Question = {
  prompt: string;
  choices: string[];
};

export type Ballot = {
  pollId: string;
  rankOrder: number[];
};
