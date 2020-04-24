# Reinforcement Learning : An introduction
这本书的中文版已经出了，我已经买了，在学校宿舍，没有带回来，非常可惜。
看看英文版也不错。

## 1 The Reinforcement Learning Problem

### 1.1 Reinforcement Learning

### 1.2 Examples

### 1.3 Elements of Reinforcement Learning

### 1.4 Limitations and Scope

### 1.5 An Extended Example: Tic-Tac-Toe

### 1.6 Summary

### 1.7 History of Reinforcement Learning

### 1.8 Bibliographical Remarks

# Tabular Solution Methods

## 2 Multi-arm Bandits
多悬臂摇臂机

### 2.1 A k-Armed Bandit Problem

### 2.2 Action-Value Methods

### 2.3 Incremental Implementation

### 2.4 Tracking a Nonstationary Problem

### 2.5 Optimistic Initial Value

### 2.6 Upper-Confidence-Bound Action Selection

### 2.7 Gradient Bandit Algorithms

### 2.8 Associative Search (Contextual Bandits)

### 2.9 Summary

## 3 Finite Markov Decision

### 3.1 The Agent-Environment Interface

### 3.2 Goals and Rewards

### 3.3 Returns 

### 3.4 Unified Notation for Episodic and Continuing Tasks

### 3.5 The Markov Property

### 3.6 Markov Decision Processes

### 3.7 Value Functions

### 3.8 Optimal Value Functions

### 3.9 Optimality and Approximation

### 3.10 Summary

## 4 Dynamic Programming

### 4.1 Policy Evaluation

### 4.2 Policy Improvement

### 4.3 Policy Iteration

### 4.4 Value Iteration

### 4.5 Asynchronous Dynamic Programming

### 4.6 Generalized Policy Iteration

### 4.7 Efficient of Dynamic Programming

### 4.8 Summary

## 5 Monte Carlo Methods

### 5.1 Monte Carlo Prediction

### 5.2 Monte Carlo Estimation of Action Values

### 5.3 Monte Carlo Control

### 5.4 Monte Carlo Control with Exploring Starts

### 5.5 Off-policy Prediction via Importance Sampling

### 5.6 Incremental Implementation

### 5.7 Off-Policy Monte Carlo Control

### 5.8 Return-Specific Importance Sampling

### 5.9 Summary

## 6 Temporal-Difference Learning

### 6.1 TD Prediction

### 6.2 Advantages of TD Prediction Methods

### 6.3 Optimality of TD(0)

### 6.4 Sarsa : On-Policy TD Control

### 6.5 Q-Learning : Off-Policy TD Control

### 6.6 Expected Sarsa

### 6.7 Maximization Bias and Double Learning

### 6.8 Games, Afterstates, and Other Special Cases

### 6.9 Summary

## 7 Multi-step Boostrapping

### 7.1 n-step TD Prediction

### 7.2 n-step Sarsa

### 7.3 n-step Off-policy Learning by Importance Sampling

### 7.4 Off-policy Learning Without Importance Sampling: The n-step Tree Backup Algorithm

### 7.5 A Unifying Algorithm: n-step Q(sigma)

### 7.6 Summary

## 8 Planning and Learning with Tabular Methods

### 8.1 Models and Planning

### 8.2 Dyna: Intergrating Planning, Acting and Learning

### 8.3 When the Model Is Wrong

### 8.4 Prioritized Sweeping

### 8.5 Planning as Part of Action Selection

### 8.6 Heuristic Search

### 8.7 Monte Carlo Tree Search

### 8.8 Summary

# 2 Approximation SOlution Methods

## 9 On-policy Prediction with Approximation

### 9.1 Value-function Approximation

### 9.2 The Prediction Objective (MSVE)

### 9.3 Stochastic-gradient and Semi-gradient Methods

### 9.4 Linear Methods

### 9.5 Feature Construction

#### 9.5.1 Polynomials

#### 9.5.2 Fourier Basis

#### 9.5.3 Coarse Coding

#### 9.5.4 Tile Coding

#### 9.5.5 Radial Basis Functions

### 9.6 Nonlinear Function Approximation: Artificial Neural Networks

### 9.7 Least-Squares TD

### 9.8 Summary

## 10 On-policy Control

### 10.1 Episodic Semi-gradient Control

### 10.2 n-step Semi-gradient Sarsa

### 10.3 Average Reward: A New Problem Setting for Continuing Tasks

### 10.4 Deprecating and Discounted Setting

### 10.5 n-step Differential Semi-gradient Sarsa

### 10.6 Summary

## 11 Off-policy Methods with Approximation

### 11.1 Semi-gradient Methods

### 11.2 Baird's Counterexample

### 11.3 The Deadly Triad

## 12 Eligibility Traces

### 12.1 The lambda-return

### 12.2 TD(lambda)

### 12.3 An On-line Forward View

### 12.4 True Online TD(lambda)

### 12.5 Dutch Traces in Monte Carlo Learning

## 13 Policy Gradient Methods

### 13.1 Policy Approximation and its Advantages

### 13.2 The Policy Gradient Theorem

### 13.3 REINFORCE: Monte Carlo Policy Gradient

### 13.4 REINFORCE with Baseline

### 13.5 Actor-Critic Methods

### 13.6 Policy Gradient for Continuing Problems (Average Reward Rate)

### 13.7 Policy Parameterization for Continuous Actions

# 3 Looking Deeper

## 14 Psychology

### 14.1 Terminology

### 14.2 Prediction and Control

### 14.3 Classical Conditioning

#### 14.3.1 The Rescorla-Wagner Model

#### 14.3.2 The TD Model

#### 14.3.3 TD Model Simulations

### 14.4 Instrumental Conditioning

### 14.5 Delayed Reinforcement

### 14.6 Cognitive Maps

### 14.7 Habitual and Goal-Directed Behaviour

### 14.8 Summary

## 15 Neuroscience

### 15.1 Neuroscience Basics

### 15.2 Reward Signals, Reinforcement Signals, Values, and Prediction Errors

### 15.3 The Reward Prediction Error Hypothesis

### 15.4 Dopamine

### 15.5 Experimental Support for the Reward Prediction Error Hypothesis

### 15.6 TD Error/ Dopamine Correspondence

### 15.7 Neural Acotr-Critic

### 15.8 Actor and Critic Learning Rules

### 15.9 Hedonistic Neurons

### 15.10 Collective Reinforcement Learning

### 15.11 Model-Based Methods in the Brain

### 15.12 Addiction

### 15.13 Summary

### 15.14 Conclusion

### 15.15 Bibliographical and Historical Remarks

## 16 Applications and Case Studies

### 16.1 TD-Gammon

### 16.2 Samuel's Checkers Player

### 16.3 The Acrobot

### 16.4 Watson's Daily-Double Wagering

### 16.5 Optimizing Memory Control

### 16.6 Human-Level Video Game Play

### 16.7 Mastering the Game of Go

### 16.8 Personalized Web Services

### 16.9 Thermal Soaring

## 17 Frontiers

### 17.1 The Unified View
