import os
from distutils.dir_util import copy_tree


def insert_code_after_string(target_string, code_to_insert):
    with open("register_commands.go", "r") as file:
        lines = file.readlines()

    for line_number, line in enumerate(lines, start=1):
        if target_string in line:
            lines.insert(line_number, code_to_insert + "\n")
            break

    with open("register_commands.go", "w") as file:
        file.writelines(lines)


command_name = input("Enter the name of the command: ")

copy_tree("commands/template", "commands/" + command_name)

os.rename(
    "commands/" + command_name + "/template.go",
    "commands/" + command_name + "/" + command_name + ".go",
)
file = open("commands/" + command_name + "/" + command_name + ".go")
contents = file.read()
replaced_contents = contents.replace("TEMPLATE", command_name)
with open("commands/" + command_name + "/" + command_name + ".go", "w") as f:
    for line in replaced_contents:
        f.write(line)
file.close()

insert_code_after_string(
    "import (", '\t"discord-go-test/commands/' + command_name + '"'
)
insert_code_after_string(
    "	commands := []*discordgo.ApplicationCommand{",
    "\t\t&" + command_name + ".CommandDefinition,",
)

insert_code_after_string(
    "	if i.Type == discordgo.InteractionApplicationCommand {",
    '\t\tif i.ApplicationCommandData().Name == "'
    + command_name
    + '"{\n\t\t\t'
    + command_name
    + ".Command(s, i)\n\t\t}",
)
