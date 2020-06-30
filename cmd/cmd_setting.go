package cmd

import (
	"io"
	"kube-recreate/pkg/util"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type CodeProperties struct {
	version, commit, branch string
}

type CmdSetting struct {
	Out                io.Writer
	Reporter           *util.Reporter
	ObjectNameProvided bool
	ObjectName         string
	Parent             *cobra.Command
	CodeProperties     CodeProperties
}

func (c *CmdSetting) Namespace() string {
	ns, err := c.Parent.Flags().GetString("namespace")
	if err != nil || len(ns) == 0 {
		namespace, _, err := genericclioptions.NewConfigFlags(true).ToRawKubeConfigLoader().Namespace()
		if err != nil || len(namespace) == 0 {
			namespace = "default"
		}
		return namespace
	}

	return ns
}

func (c *CmdSetting) AllFlag() bool {
	all, err := c.Parent.Flags().GetBool("all")
	if err != nil {
		return false
	}
	return all
}

func (c *CmdSetting) AllNamespacesFlag() bool {
	all, err := c.Parent.Flags().GetBool("all-namespaces")
	if err != nil {
		return false
	}
	return all
}

type CmdSettingOption func(*CmdSetting)

func NewCmdSettings(opts ...CmdSettingOption) *CmdSetting {
	var (
		defaultOut                = os.Stdout
		defaultReporter           = util.NewReporter(defaultOut)
		defaultObjectNameProvided = false
		defaultObjectName         = ""
		defaultCodeProperties     = CodeProperties{commit: "dev", branch: "dev", version: "dev"}
	)

	h := &CmdSetting{
		Out:                defaultOut,
		Reporter:           defaultReporter,
		ObjectNameProvided: defaultObjectNameProvided,
		ObjectName:         defaultObjectName,
		CodeProperties:     defaultCodeProperties,
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

func WithObjectName(objectName string) CmdSettingOption {
	return func(h *CmdSetting) {
		h.ObjectName = objectName
		h.ObjectNameProvided = true
	}
}

func WithOutWriter(out io.Writer) CmdSettingOption {
	return func(h *CmdSetting) {
		h.Out = out
	}
}

func WithParentCmd(parent *cobra.Command) CmdSettingOption {
	return func(h *CmdSetting) {
		h.Parent = parent
	}
}

func WithCodeVersion(version string) CmdSettingOption {
	return func(h *CmdSetting) {
		h.CodeProperties.version = version
	}
}

func WithCodeCommit(commit string) CmdSettingOption {
	return func(h *CmdSetting) {
		h.CodeProperties.commit = commit
	}
}

func WithCodeBranch(branch string) CmdSettingOption {
	return func(h *CmdSetting) {
		h.CodeProperties.branch = branch
	}
}
