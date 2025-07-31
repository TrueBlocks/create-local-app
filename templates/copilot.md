# COPILOT BEHAVIORAL REQUIREMENTS - STRICT COMPLIANCE MANDATORY

## ABSOLUTE PROHIBITIONS - ZERO TOLERANCE

**❌ NEVER USE NPM OR NPX** - YARN ONLY, NO EXCEPTIONS
- If you type `npm` or `npx` even once, you have failed
- ALWAYS use `yarn` for ALL package management
- ALWAYS run `yarn start` from repo root, NEVER from frontend folder
- ALWAYS run `yarn build` from repo root, NEVER from frontend folder
- ALWAYS run `yarn test` from repo root if testing required

**❌ NEVER USE console.log IN FRONTEND CODE**
- ALWAYS use `Log from @utils` - single string parameter only
- This outputs to command line where I can see it
- console.log is invisible in Wails apps - you're wasting both our time

**❌ NEVER RUN APP IN BROWSER MODE**
- NEVER use localhost:5173 or browser preview
- ALWAYS use `yarn start` from root (Wails dev mode)
- Browser mode doesn't work properly with this Wails app

**❌ NEVER MAKE FALSE CLAIMS ABOUT COMPLETENESS**
- Don't claim "full implementation" when you built mocks
- Don't claim "working transactions" when you generate fake hashes
- Don't claim "backend integration" when data is hardcoded
- BE BRUTALLY HONEST about what's real vs. mocked

## MANDATORY BEHAVIORS - STRICT ADHERENCE REQUIRED

**✅ ALWAYS CHECK FILE CONTENTS BEFORE EDITING**
- Read current file state with read_file tool FIRST
- User makes manual edits between requests - files change
- Assuming file contents leads to merge conflicts and wasted time

**✅ ALWAYS USE EXISTING PATTERNS AND COMPONENTS**
- Don't reinvent Table components - use existing BaseTab/Table patterns
- Follow established naming conventions from other collections
- Reuse existing hooks, utilities, and styling patterns
- Study existing code before creating new patterns

**✅ ALWAYS FOLLOW ESTABLISHED ARCHITECTURE**
- Use DataFacet enum values, don't make up custom strings
- Follow Collection/Store/Page patterns from other views
- Use existing type definitions from @models
- Integrate with existing routing and state management

**✅ ALWAYS VERIFY AND VALIDATE**
- Search codebase for existing implementations before creating new ones
- Check for existing utilities before writing duplicates
- Validate TypeScript types against existing models
- Test integration points with existing systems

## WORK PROCESS DEMANDS

**📋 TODO LIST COMPLIANCE**
- Work through ToDoList systematically, one step at a time
- Mark ✅ ONLY when step is actually complete (not just coded)
- STOP at Checkpoints and wait for my confirmation
- Move completed tasks to "Completed Tasks" section only after my approval

**🔍 ANALYSIS REQUIREMENTS**
- When asked "Where are we?" provide honest status assessment
- Clearly distinguish between completed, in-progress, and not-started work
- Identify gaps between claimed completion and actual functionality
- Be specific about what still needs real implementation vs. mocking

**❓ CLARIFICATION PROTOCOL**
- Ask specific questions rather than making assumptions
- Request clarification with: "Please clarify: [specific question]"
- Don't proceed with guesswork when requirements are unclear
- Acknowledge when you don't understand something

## ERROR HANDLING DEMANDS

**🛑 IMMEDIATE STOP CONDITIONS**
- If tests fail: STOP, report exact error, wait for instructions
- If linting fails: STOP, report issues, wait for fixes
- If build fails: STOP, provide full error output, wait for guidance
- If unsure about requirements: STOP, ask for clarification

**🚫 WHAT NOT TO DO WHEN ERRORS OCCUR**
- Don't try to "fix" test failures yourself
- Don't ignore linting warnings
- Don't work around build issues
- Don't guess at error solutions

## CODING STANDARDS - NON-NEGOTIABLE

**TypeScript Requirements:**
- Never import React explicitly (implicitly available)
- Use proper types from @models, don't create duplicates
- Follow existing interfaces and type patterns
- Use defensive coding with proper error boundaries

**Styling Requirements:**
- Use Mantine components and styling system
- Avoid inline styles - create .css files for custom styles
- Follow existing color schemes and spacing patterns
- Maintain responsive design consistency

**Code Organization:**
- No comments in production code
- Follow existing folder structure patterns
- Use established import/export conventions
- Maintain separation of concerns

## FINAL WARNING

Violating these requirements wastes my time and yours. I've given you these instructions multiple times throughout our conversation. Follow them exactly, or expect to be corrected and have to redo work.

The goal is efficient, high-quality development that integrates seamlessly with existing codebase patterns. Deviation from these requirements achieves the opposite.