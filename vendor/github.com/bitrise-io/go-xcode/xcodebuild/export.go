package xcodebuild

import (
	"bytes"
	"io"
	"os"
	"os/exec"

	"github.com/bitrise-io/go-utils/command"
)

/*
xcodebuild -exportArchive \
	-archivePath <xcarchivepath> \
	-exportPath <destinationpath> \
	-exportOptionsPlist <plistpath>
*/

// ExportCommandModel ...
type ExportCommandModel struct {
	archivePath        string
	exportDir          string
	exportOptionsPlist string
	authentication     *AuthenticationParams
}

// NewExportCommand ...
func NewExportCommand() *ExportCommandModel {
	return &ExportCommandModel{}
}

// SetArchivePath ...
func (c *ExportCommandModel) SetArchivePath(archivePath string) *ExportCommandModel {
	c.archivePath = archivePath
	return c
}

// SetExportDir ...
func (c *ExportCommandModel) SetExportDir(exportDir string) *ExportCommandModel {
	c.exportDir = exportDir
	return c
}

// SetExportOptionsPlist ...
func (c *ExportCommandModel) SetExportOptionsPlist(exportOptionsPlist string) *ExportCommandModel {
	c.exportOptionsPlist = exportOptionsPlist
	return c
}

// SetAuthentication ...
func (c *ExportCommandModel) SetAuthentication(authenticationParams AuthenticationParams) *ExportCommandModel {
	c.authentication = &authenticationParams
	return c
}

func (c *ExportCommandModel) cmdSlice() []string {
	slice := []string{toolName}
	slice = append(slice, c.CommandArgs()...)

	return slice
}

// CommandArgs returns the xcodebuild command arguments for the export action
func (c *ExportCommandModel) CommandArgs() []string {
	slice := []string{"-exportArchive"}
	if c.archivePath != "" {
		slice = append(slice, "-archivePath", c.archivePath)
	}

	if c.exportDir != "" {
		slice = append(slice, "-exportPath", c.exportDir)
	}

	if c.exportOptionsPlist != "" {
		slice = append(slice, "-exportOptionsPlist", c.exportOptionsPlist)
	}

	if c.authentication != nil {
		slice = append(slice, c.authentication.args()...)
	}

	return slice
}

// PrintableCmd ...
func (c *ExportCommandModel) PrintableCmd() string {
	cmdSlice := c.cmdSlice()
	return command.PrintableCommandArgs(false, cmdSlice)
}

// Command ...
func (c *ExportCommandModel) Command() *command.Model {
	cmdSlice := c.cmdSlice()
	return command.New(cmdSlice[0], cmdSlice[1:]...)
}

// Cmd ...
func (c *ExportCommandModel) Cmd() *exec.Cmd {
	command := c.Command()
	return command.GetCmd()
}

// Run ...
func (c *ExportCommandModel) Run() error {
	command := c.Command()

	command.SetStdout(os.Stdout)
	command.SetStderr(os.Stderr)

	return command.Run()
}

// RunAndReturnOutput ...
func (c *ExportCommandModel) RunAndReturnOutput() (string, error) {
	command := c.Command()

	var outBuffer bytes.Buffer
	outWriter := io.MultiWriter(&outBuffer, os.Stdout)

	command.SetStdout(outWriter)
	command.SetStderr(outWriter)

	err := command.Run()
	out := outBuffer.String()

	return out, err
}
