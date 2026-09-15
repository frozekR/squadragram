package model

type (
	CharacterRole string
	SkillType     string
	OwnerType     string
	MediaType     string
)

const (
	CharacterRoleDamage    = CharacterRole("DAMAGE")
	CharacterRoleTank      = CharacterRole("TANK")
	CharacterRoleTechnical = CharacterRole("TECHNICAL")
)

const (
	SkillTypePassive     = SkillType("PASSIVE")
	SkillTypeRush        = SkillType("RUSH_ATTACK")
	SkillTypeSkill       = SkillType("SKILL")
	SkillTypeSuperAttack = SkillType("SUPER_ATTACK")
	SkillTypeTransform   = SkillType("TRANSFORMATION")
)

const (
	OwnerTypeCharacter = OwnerType("CHARACTER")
	OwnerTypeSkill     = OwnerType("SKILL")
)

const (
	MediaTypeIcon   = MediaType("ICON")
	MediaTypeRender = MediaType("RENDER")
	MediaTypeDemo   = MediaType("DEMO")
)
