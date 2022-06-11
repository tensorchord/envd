package types

const (
	ContainerLabelName        = "ai.tensorchord.envd.name"
	ContainerLabelJupyterAddr = "ai.tensorchord.envd.jupyter.address"
	ContainerLabelSSHPort     = "ai.tensorchord.envd.ssh.port"

	ImageLabelVendor  = "ai.tensorchord.envd.vendor"
	ImageLabelGPU     = "ai.tensorchord.envd.gpu"
	ImageLabelAPT     = "ai.tensorchord.envd.apt.packages"
	ImageLabelPyPI    = "ai.tensorchord.envd.pypi.packages"
	ImageLabelR       = "ai.tensorchord.envd.r.packages"
	ImageLabelCUDA    = "ai.tensorchord.envd.gpu.cuda"
	ImageLabelCUDNN   = "ai.tensorchord.envd.gpu.cudnn"
	ImageLabelContext = "ai.tensorchord.envd.build.context"

	ImageVendorEnvd = "envd"
)
