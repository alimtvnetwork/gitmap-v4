import { Link } from "react-router-dom";
import { FolderGit2, GitBranch, RefreshCw, Eye } from "lucide-react";
import DocsLayout from "@/components/docs/DocsLayout";
import FeatureCard from "@/components/docs/FeatureCard";
import InstallBlock from "@/components/docs/InstallBlock";
import { VERSION } from "@/constants/index";

const HomePage = () => {
  return (
    <DocsLayout>
      <section className="py-12 text-center">
        <div className="flex items-center justify-center gap-3 mb-4">
          <h1 className="text-4xl md:text-5xl font-heading font-bold docs-h1">
            gitmap
          </h1>
          <span className="px-2 py-0.5 rounded text-xs font-mono bg-primary/10 text-primary border border-primary/20">
            {VERSION}
          </span>
        </div>
        <p className="text-lg text-muted-foreground max-w-2xl mx-auto mb-8 leading-relaxed font-sans">
          Scan a folder tree for Git repos, generate structured clone files, and
          re-clone the exact layout on any machine. Track, group, release, and
          manage repositories from a single CLI.
        </p>
        <div className="mb-8 max-w-3xl mx-auto">
          <InstallBlock
            tabs={[
              {
                label: "Windows",
                command:
                  "irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v3/main/install-quick.ps1 | iex",
              },
              {
                label: "Linux / macOS",
                command:
                  "curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v3/main/install-quick.sh | bash",
              },
            ]}
          />
        </div>
        <div className="flex gap-4 justify-center">
          <Link
            to="/getting-started"
            className="px-5 py-2.5 rounded-lg bg-primary text-primary-foreground font-mono text-sm font-medium hover:bg-primary/90 transition-colors"
          >
            Get Started
          </Link>
          <Link
            to="/commands"
            className="px-5 py-2.5 rounded-lg border border-border text-foreground font-mono text-sm font-medium hover:bg-muted transition-colors"
          >
            View Commands
          </Link>
        </div>
      </section>

      <hr className="docs-hr" />

      <section className="grid md:grid-cols-2 gap-4 py-8">
        <FeatureCard
          icon={FolderGit2}
          title="Scan & Map"
          description="Recursively discover Git repos, extract metadata, and output CSV/JSON/terminal views with clone scripts."
        />
        <FeatureCard
          icon={GitBranch}
          title="Clone & Restore"
          description="Re-clone the exact folder structure on a new machine from JSON, CSV, or text files with safe-pull and progress tracking."
        />
        <FeatureCard
          icon={RefreshCw}
          title="Release & Version"
          description="Create releases with tags, branches, changelogs, and semantic versioning — all from the command line."
        />
        <FeatureCard
          icon={Eye}
          title="Watch & Monitor"
          description="Live-refresh dashboard showing dirty/clean status, ahead/behind counts, and stash entries across all tracked repos."
        />
      </section>

      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{
          __html: JSON.stringify({
            "@context": "https://schema.org",
            "@type": "SoftwareApplication",
            name: "gitmap",
            applicationCategory: "DeveloperApplication",
            operatingSystem: "Windows, macOS, Linux",
            description: "CLI tool to scan, map, and re-clone Git repository trees.",
          }),
        }}
      />
    </DocsLayout>
  );
};

export default HomePage;
