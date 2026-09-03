// package main
//
// import (
// 	"fmt"
// 	"strings"
// )
//
// type Command struct {
// 	Name    string
// 	MinArgs int
// 	MaxArgs int
// 	// Description	string
// 	Validator func([]string) error
// 	Handler   func(*Player, []string) error
// 	// NOTE: maybe change the design
// 	RequiresAuth bool
// }
//
// // registry with all the command available
// type CommandRegistry struct {
// 	commands map[string]*Command
// }
//
// // command registry constructor
// func NewCommandRegistry() *CommandRegistry {
// 	cr := &CommandRegistry{
// 		commands: make(map[string]*Command),
// 	}
// 	cr.registerCommands()
// 	return cr
// }
//
// func (cr *CommandRegistry) registerCommands() {
// 	// LOOK command
// 	cr.commands["LOOK"] = &Command{
// 		Name:         "LOOK",
// 		MinArgs:      0,
// 		MaxArgs:      0,
// 		RequiresAuth: true,
// 		Validator: func(args []string) error {
// 			if len(args) > 0 {
// 				// NOTE: logg
// 				return fmt.Errorf("LOOK takes no arguments")
// 			}
// 			return nil
// 		},
// 		Handler: func(p *Player, args []string) error {
// 			// place holder for look logic
// 			return nil
// 		},
// 	}
// 	cr.commands["MOVE"] = &Command{
// 		// place holder
// 	}
// 	// and so on
// }
//
// func (s *Server) handleCommand(player *Player, line string) {
// 	parts := strings.Fields(line)
// 	if len(parts) == 0 {
// 		return
// 	}
// 	// NOTE: later we need to change this to forbid lowercase COMMAND
// 	commandName := strings.ToUpper(parts[0])
// 	args := parts[1:]
// 	cmd, exists := s.cmdRegistry.commands[commandName]
// 	if !exists {
// 		// NOTE: add logg -- maybe custom error code on unknown command
// 		s.sendError(player, 400, fmt.Sprintf("UNKNOWN_COMMAND: %s", commandName))
// 		return
// 	}
// 	if cmd.RequiresAuth && player.State != Authenticated {
// 		s.sendError(player, 401, "NOT_AUTHENTICATED")
// 	}
// 	if len(args) < cmd.MinArgs {
// 		s.sendError(player, 400, fmt.Sprintf("TOO_FEW_ARGS: Need at least %d arguments", cmd.MinArgs))
// 		return
// 	}
// 	if cmd.MaxArgs > 0 && len(args) > cmd.MaxArgs {
// 		s.sendError(player, 400, fmt.Sprintf("TOO_MANY_ARGS: Maximum %d arguments allowed", cmd.MaxArgs))
// 		return
// 	}
// 	if cmd.Validator != nil {
// 		if err := cmd.Validator(args); err != nil {
// 			s.sendError(player, 400, fmt.Sprintf("INVALID_ARGS: %s", err.Error()))
// 			return
// 		}
// 	}
// 	if err := cmd.Handler(player, args); err != nil {
// 		s.sendError(player, 500, fmt.Sprintf("COMMAND_ERROR: %s", err.Error()))
// 		return
// 	}
// }
