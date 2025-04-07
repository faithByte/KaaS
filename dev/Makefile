MAKEFLAGS += --silent --ignore-errors

PLAYBOOK = playbook.yml

## vagrant
vms:
	vagrant up
	vagrant ssh master -c "sudo -i"

master:
	vagrant up master
	vagrant ssh master -c "sudo -i"

clean:
	vagrant destroy -f

re: clean vms

## ansible
# all:
# 	@echo $(PLAYBOOK)

check:
	@cd ../project && \
	ansible-playbook $(PLAYBOOK) --syntax-check && \
	ansible-playbook $(PLAYBOOK) --check

run:
	@cd ../project && ansible-playbook $(PLAYBOOK);

slurm:
	@cd ../slurm && ansible-playbook $(PLAYBOOK);

# task:
# 	@cd ../project && ansible-playbook $(PLAYBOOK) --start-at-task="$(TASK)"

.PHONY: all vms master clean re run check 