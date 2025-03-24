# Verdict

Verdict is a full-stack web app for creating ranked choice polls, in which voters rank each choice by preference instead of choosing only one. The instant runoff algorithm calculates winners by repeatedly eliminating the choice with the fewest votes, redistributing ballots until a single choice achieves a strict majority.

## Table of Contents

1. [Tech Stack](#tech-stack)
2. [Current Features](#current-features)
3. [About Ranked Choice Voting](#about-ranked-choice-voting)

## Tech Stack 

| | |
| - | - |
| ***Backend*** | ![Static Badge](https://img.shields.io/badge/Go-00ADD8) ![Static Badge](https://img.shields.io/badge/Lambda-FF9900) ![Static Badge](https://img.shields.io/badge/API_Gateway-FF4F8B) ![Static Badge](https://img.shields.io/badge/DynamoDB-4053D6) |
| ***Frontend*** | ![Static Badge](https://img.shields.io/badge/TypeScript-3178C6) ![Static Badge](https://img.shields.io/badge/React-61DAFB) ![Static Badge](https://img.shields.io/badge/Vite-646CFF) ![Static Badge](https://img.shields.io/badge/Vitest-6E9F18) ![Static Badge](https://img.shields.io/badge/pnpm-F69220) |
| ***Dev Tools*** | ![Static Badge](https://img.shields.io/badge/Docker-2496ED) |

## Current Features

- Create polls
- Cast ballots
- Calculate results

## About Ranked Choice Voting

Verdict implements a voting algorithm known as ranked choice voting or instant runoff voting. There are several variations, but this is the one used here:

### Voting process

Instead of selecting a single choice, voters rank each choice as their first choice, second choice, third choice, etc., all the way to last choice.

### Determining a winner

Instead of immediately selecting a choice that has only a plurality of votes, the algorithm first checks if any choice has a strict majority of votes. If no choice does, the choice with the fewest votes is eliminated, and its votes are redistributed to the voters' next highest choices. This process of elimination continues until a single choice has a strict majority of votes.

### What about ties for last?

In each round, the choice with the fewest votes is eliminated, but what if multiple choices are tied for last place? In this case, a sub-poll is simulated between only the tied choices. This is possible because voters provide a rank for every choice, allowing the algorithm to determine their preferences amongst any subset of choices. If there is another tie for last place, the tie-breaking algorithm continues recursively.

While unlikely unless the numbers of choices and voters are very small, it is possible that multiple choices are tied for last place and received perfectly equivalent rankings. In this case, one of these lowest ranking choices is eliminated by pseudorandom number generation.
