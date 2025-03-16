export type Question = {
  prompt: string;
  choices: string[];
};

export type Ballot = {
  pollId: string;
  rankOrder: number[];
};

export type Result = {
  prompt: string;
  winningChoice: string;
  numVotes: number;
  winningRound: number;
};
