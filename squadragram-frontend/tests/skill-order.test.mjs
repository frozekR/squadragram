import assert from "node:assert/strict";
import test from "node:test";
import { groupSkills } from "../src/shared/types/skill.ts";

test("categories follow the moveset order, positions sort within each category", () => {
  const skills = [
    { id: 4, type: "SKILL", sort_order: 20 },
    { id: 5, type: "SUPER_ATTACK", sort_order: 0 },
    { id: 3, type: "SKILL", sort_order: 5 },
    { id: 1, type: "PASSIVE", sort_order: 100 },
    { id: 2, type: "SKILL", sort_order: 5 },
  ];
  const original = structuredClone(skills);
  const groups = groupSkills(skills);
  assert.deepEqual(groups.map(group => group.type), ["PASSIVE", "SKILL", "SUPER_ATTACK"]);
  assert.deepEqual(groups[1].skills.map(skill => skill.id), [2, 3, 4]);
  assert.deepEqual(skills, original, "sorting must not mutate the query cache");
});

test("empty categories are omitted and old records default to position zero", () => {
  assert.deepEqual(groupSkills([]), []);
  const groups = groupSkills([{ id: 2, type: "SKILL", sort_order: 1 }, { id: 1, type: "SKILL" }]);
  assert.deepEqual(groups[0].skills.map(skill => skill.id), [1, 2]);
});
