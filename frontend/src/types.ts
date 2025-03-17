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
  totalVotes: number;
  winningVotes: number;
  winningChoice: string;
  winningRound: number;
};
