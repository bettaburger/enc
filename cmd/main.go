package main

import (
	"enc/internal"
	"fmt"
	"os"
	
	tea "charm.land/bubbletea/v2"
)

type model struct {
	options 		[]	string		// main menu options
	cursor 				int 		// what our cursor is pointing to
	selected 			int			// what option is being selected
}

func initalModel() model {
	return model {
		options: 	[]string{"url", "help", "exit"},
		selected: 	-1	, // nothing selected
		cursor: 	0, 
	}
}

func (m model) Init() tea.Cmd {
    return nil
}

// KEYBOARD CONTROL
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyPressMsg:
        switch msg.String() {
        // keys to exit the process
        case "ctrl+c", "q":
            return m, tea.Quit

        // move the cursor up
        case "up":
            if m.cursor > 0 {
                m.cursor--
            }
        // move the cursor down
        case "down":
            if m.cursor < len(m.options)-1 {
                m.cursor++
            }
        // selected state
        case "enter", "space":
            if m.selected == m.cursor {
				m.selected = -1
                
            } else {
                m.selected = m.cursor
            }
        }
    }
    return m, nil
}

// RENDERS UI 
func (m model) View() tea.View {
    s := "Tool for http look up:\n\n"
    for i, option := range m.options {
        cursor := " " 
        if m.cursor == i {
            cursor = ">" // cursor
        }
        checked := " " 
        if m.selected == i {
            checked = "x" // selected
        }
        // Render the row
        s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, option)
    }
    s += "\nPress q to quit.\n"
    return tea.NewView(s)
}


func main() {
	program := tea.NewProgram(initalModel())
	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Program error", err)
	}

	internal.Fetch()
	// run tool here
}

