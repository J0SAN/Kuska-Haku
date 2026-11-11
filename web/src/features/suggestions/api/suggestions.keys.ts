export const suggestionsKeys = {
  all: ["suggestions"] as const,
  list: (onlyNew: boolean) =>
    [...suggestionsKeys.all, "list", onlyNew] as const,
};
