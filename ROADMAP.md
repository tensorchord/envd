# envd Roadmap

Want to jump in? Welcome discussions and contributions! 

- Chat with us on [💬 Discord](https://discord.gg/KqswhpVgdU)
- Have a look at [`good first issue 💖`](https://github.com/tensorchord/envd/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue+%E2%9D%A4%EF%B8%8F%22) issues!
- More on [contributing page](https://envd.tensorchord.ai/docs/community/contributing)

## Short term - what we're working on now 🎉

We are working on building the MVP!

- **Build language**
    - Support more languages. See [feat: Support environments with multiple languages](https://github.com/tensorchord/envd/issues/407), [feat: Support Julia language](https://github.com/tensorchord/envd/issues/408)
    - [Support multi-target builds](https://github.com/tensorchord/envd/issues/403)
    - [Implement our own language server](https://github.com/tensorchord/envd/issues/358)
    - [Support oh-my-zsh plugins configuration](https://github.com/tensorchord/envd/issues/106)
    - [Customize the base images](https://github.com/tensorchord/envd/issues/261)
- **Runtime**
    - [Add agent to collect metrics in the container (or in the host)](https://github.com/tensorchord/envd/issues/218). We need to collect and show the metrics of GPUs/CPUs to users, maybe via a web-based UI.
    - [Support buildkitd moby worker in dockerd 22.06-beta](https://github.com/tensorchord/envd/issues/51). It is a huge enhancement which accelarates the build process. We use moby worker in docker 22.06, and fallback to docker worker in docker 20.10.
- **Ecosystem**
    - [Contribute the vscode-envd extension](https://github.com/tensorchord/vscode-envd)
    - Provide web-based UI

## Medium term - what we're working on next! 🏃

- **Language**
    - [Support data management in the environment](https://github.com/tensorchord/envd/issues/5)
- **Runtime**
    - [Support both Kubernetes and Docker runtimes](https://github.com/tensorchord/envd/issues/179). Data science or AI/ML teams usually do not only develop on their laptops. They need GPU clusters to run the large scale training jobs. envd supports both Kubernetes and Docker runtimes. Thus users can move to cloud without any code change.
- **Ecosystem**
    - [Design the extension mechanism to reuse user-defined build funcs](https://github.com/tensorchord/envd/issues/91). Users can run `load(<custom-package>)` to load and use their own build functions.

## Longer term items - working on this soon! ⏩

These items are not essiential for an MVP, but are a part of our longer term plans. Feel free to jump in on these if you're interested!

- **Runtime**
    - Continuous profiler that continuously collects line-level profiling performance data from envd environments, helps engineers find bottlenecks in the training code.
    - [Support the OCI runtime spec-compatible runtime](https://github.com/tensorchord/envd/issues/282)
- **Ecosystem**
    - Integrate with other open source development tools e.g. [aim](https://github.com/aimhubio/aim).

To be added in the future.
