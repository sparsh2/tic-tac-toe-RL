# tic-tac-toe-RL
This is a simple demonstration of how an RL agent can learn to play tic-tac-toe by playing against itself and updating its knowledge table using Temporal-Difference method.

This is inspired from the example provided in the introduction of the book [Reinforcement Learning, second edition: An Introduction](https://www.amazon.com/dp/0262039249?ref_=cm_sw_r_cp_ud_dp_PXA1NE9WXG9ZJGSRM7TY)

# How to play against the agent?

Run `go run main.go` from root. Additionally set number of training episodes in [main.go](./main.go#L99) to get the desired strength of the agent. By 10,000 training episodes, the agent plays optimally.
