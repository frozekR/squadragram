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
	CharacterRoleMelee     = CharacterRole("MELEE")
	CharacterRoleRanged    = CharacterRole("RANGED")
)

const (
	SkillTypePassive        = SkillType("PASSIVE")
	SkillTypeRush           = SkillType("RUSH_ATTACK")
	SkillTypeSkill          = SkillType("SKILL")
	SkillTypeSuperAttack    = SkillType("SUPER_ATTACK")
	SkillTypeMaxSuperAttack = SkillType("MAX_SUPER_ATTACK")
	SkillTypeTransform      = SkillType("TRANSFORMATION")
)

const (
	OwnerTypeCharacter = OwnerType("CHARACTER")
	OwnerTypeSkill     = OwnerType("SKILL")
	OwnerTypeEmote     = OwnerType("EMOTE")
	OwnerTypeSkin      = OwnerType("SKIN")
)

const (
	MediaTypeIcon   = MediaType("ICON")
	MediaTypeRender = MediaType("RENDER")
	MediaTypeDemo   = MediaType("DEMO")
)
