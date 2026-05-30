import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { DataTable, type Column } from "./DataTable";

interface TestRow {
  id: number;
  name: string;
  level: number;
  isActive: boolean;
}

const baseColumns: Column<TestRow>[] = [
  { header: "Name", accessor: "name" },
  { header: "Level", accessor: "level", align: "center" as const },
  { header: "Active", accessor: "isActive" },
];

const baseData: TestRow[] = [
  { id: 1, name: "Fireball", level: 5, isActive: true },
  { id: 2, name: "Icebolt", level: 3, isActive: false },
];

describe("DataTable", () => {
  describe("table element", () => {
    it("renders a table element", () => {
      render(<DataTable columns={baseColumns} data={baseData} getKey={(row) => row.id} />);

      expect(document.querySelector("table")).toBeInTheDocument();
    });

    it("renders table headers", () => {
      render(<DataTable columns={baseColumns} data={baseData} getKey={(row) => row.id} />);

      const headers = document.querySelectorAll("thead th");
      expect(headers).toHaveLength(3);
      expect(headers[0]).toHaveTextContent("Name");
      expect(headers[1]).toHaveTextContent("Level");
      expect(headers[2]).toHaveTextContent("Active");
    });

    it("renders table rows", () => {
      render(<DataTable columns={baseColumns} data={baseData} getKey={(row) => row.id} />);

      const rows = document.querySelectorAll("tbody tr");
      expect(rows).toHaveLength(2);
    });

    it("renders cell data", () => {
      render(<DataTable columns={baseColumns} data={baseData} getKey={(row) => row.id} />);

      const cells = document.querySelectorAll("tbody td");
      const cellTexts = Array.from(cells).map((c) => c.textContent);
      expect(cellTexts).toContain("Fireball");
      expect(cellTexts).toContain("Icebolt");
      expect(cellTexts).toContain("5");
      expect(cellTexts).toContain("3");
    });
  });

  describe("empty state", () => {
    it("renders empty message when data is empty", () => {
      render(
        <DataTable
          columns={baseColumns}
          data={[]}
          getKey={(row) => row.id}
          emptyMessage="No abilities found."
        />,
      );

      expect(screen.getByText("No abilities found.")).toBeInTheDocument();
      expect(document.querySelector("table")).not.toBeInTheDocument();
    });

    it("uses default empty message", () => {
      render(<DataTable columns={baseColumns} data={[]} getKey={(row) => row.id} />);

      expect(screen.getByText("No records found.")).toBeInTheDocument();
    });
  });

  describe("boolean rendering", () => {
    it("renders badge for true boolean values", () => {
      render(<DataTable columns={baseColumns} data={baseData} getKey={(row) => row.id} />);

      const badges = document.querySelectorAll(".badge");
      const successBadges = Array.from(badges).filter(b => b.classList.contains("badge-success"));
      expect(successBadges.length).toBeGreaterThan(0);
    });

    it("renders badge for false boolean values", () => {
      render(<DataTable columns={baseColumns} data={baseData} getKey={(row) => row.id} />);

      const badges = document.querySelectorAll(".badge");
      const neutralBadges = Array.from(badges).filter(b => b.classList.contains("badge-neutral"));
      expect(neutralBadges.length).toBeGreaterThan(0);
    });
  });

  describe("null/undefined handling", () => {
    it("renders em dash for null values", () => {
      const dataWithNull: TestRow[] = [{ id: 1, name: null as unknown as string, level: 1, isActive: true }];
      render(<DataTable columns={baseColumns} data={dataWithNull} getKey={(row) => row.id} />);

      const muted = document.querySelectorAll(".text-muted");
      const emDash = Array.from(muted).find(e => e.textContent === "—");
      expect(emDash).toBeInTheDocument();
    });
  });

  describe("custom render function", () => {
    it("uses custom render for column", () => {
      const columns: Column<TestRow>[] = [
        { header: "Name", accessor: "name", render: (val) => `**${val}**` },
        { header: "Level", accessor: "level" },
      ];

      render(<DataTable columns={columns} data={baseData} getKey={(row) => row.id} />);

      const cells = document.querySelectorAll("tbody td");
      expect(cells[0]).toHaveTextContent("**Fireball**");
    });
  });

  describe("row click", () => {
    it("calls onRowClick with row data", () => {
      const handleClick = vi.fn();
      render(
        <DataTable columns={baseColumns} data={baseData} getKey={(row) => row.id} onRowClick={handleClick} />,
      );

      const firstRow = document.querySelector("tbody tr");
      firstRow?.click();

      expect(handleClick).toHaveBeenCalledWith(baseData[0]);
    });

    it("row is clickable when onRowClick provided", () => {
      render(
        <DataTable columns={baseColumns} data={baseData} getKey={(row) => row.id} onRowClick={vi.fn()} />,
      );

      const firstRow = document.querySelector("tbody tr");
      expect(firstRow).toHaveClass("clickable-row");
    });
  });

  describe("expanded row", () => {
    it("renders expanded row content when prop is provided", () => {
      const expandedRow = (row: TestRow) => <div data-testid="expanded">Details for {row.name}</div>;
      render(
        <DataTable
          columns={baseColumns}
          data={baseData}
          getKey={(row) => row.id}
          expandedRow={expandedRow}
        />,
      );

      // There are 4 because DataTable renders both mobile and desktop views
      const expandedElements = screen.getAllByTestId("expanded");
      expect(expandedElements.length).toBeGreaterThanOrEqual(2);
    });

    it("does not render expanded rows when prop not provided", () => {
      render(<DataTable columns={baseColumns} data={baseData} getKey={(row) => row.id} />);

      expect(screen.queryByTestId("expanded")).not.toBeInTheDocument();
    });
  });

  describe("alignment", () => {
    it("applies center alignment class to header", () => {
      render(<DataTable columns={baseColumns} data={baseData} getKey={(row) => row.id} />);

      const levelHeader = document.querySelectorAll("thead th")[1];
      expect(levelHeader).toHaveClass("text-center");
    });

    it("applies center alignment to cells", () => {
      render(<DataTable columns={baseColumns} data={baseData} getKey={(row) => row.id} />);

      const levelCell = document.querySelectorAll("tbody td")[1];
      expect(levelCell).toHaveClass("text-center");
    });
  });

  describe("accessor path resolution", () => {
    it("resolves dot-notation accessors", () => {
      type NestedRow = { id: number; stats: { level: number; damage: number } };
      const columns: Column<NestedRow>[] = [
        { header: "Level", accessor: "stats.level" },
      ];
      const data: NestedRow[] = [{ id: 1, stats: { level: 5, damage: 10 } }];

      render(<DataTable columns={columns} data={data} getKey={(row) => row.id} />);

      expect(document.querySelector("tbody td")).toHaveTextContent("5");
    });

    it("returns undefined em dash for missing nested path", () => {
      type PartialRow = { id: number; stats?: { level: number } };
      const columns: Column<PartialRow>[] = [
        { header: "Level", accessor: "stats.level" },
      ];
      const data: PartialRow[] = [{ id: 1 }];

      render(<DataTable columns={columns} data={data} getKey={(row) => row.id} />);

      expect(document.querySelector(".text-muted")).toHaveTextContent("—");
    });
  });

  describe("className prop", () => {
    it("applies custom className to wrapper", () => {
      render(
        <DataTable columns={baseColumns} data={baseData} getKey={(row) => row.id} className="custom-table" />,
      );

      expect(document.querySelector(".custom-table")).toBeInTheDocument();
    });
  });

  describe("variant prop", () => {
    it("applies table-dark class when variant is dark", () => {
      render(
        <DataTable columns={baseColumns} data={baseData} getKey={(row) => row.id} variant="dark" />,
      );

      expect(document.querySelector(".table-dark")).toBeInTheDocument();
    });

    it("applies table class when variant is default", () => {
      render(
        <DataTable columns={baseColumns} data={baseData} getKey={(row) => row.id} variant="default" />,
      );

      expect(document.querySelector(".table")).toBeInTheDocument();
    });
  });
});
